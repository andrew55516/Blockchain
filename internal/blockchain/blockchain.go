package blockchain

import (
	"Blockchain/cmd/utils"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
)

const (
	MINING_DIFFICULTY = 6
	MINING_SENDER     = "THE BLOCKCHAIN"
	MINING_REWARD     = 1.0
)

type Blockchain struct {
	transactionPool   []*Transaction
	chain             []*Block
	blockchainAddress string
	port              uint16
}

func NewBlockchain(blockchainAdress string, port uint16) *Blockchain {
	b := &Block{}
	bc := new(Blockchain)
	bc.blockchainAddress = blockchainAdress
	bc.CreateBlock(0, b.Hash())
	bc.port = port
	return bc
}

func (bc *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"chain"`
	}{Blocks: bc.chain})
}

func (bc *Blockchain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := NewBlock(nonce, previousHash, bc.transactionPool)
	bc.chain = append(bc.chain, b)
	bc.transactionPool = []*Transaction{}
	return b
}

func (bc *Blockchain) Print() {
	for i, block := range bc.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
		fmt.Println()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 80))
}

func (bc *Blockchain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) AddTransaction(sender string, recipient string, value float32,
	senderPublicKey *ecdsa.PublicKey, s *utils.Signature) bool {
	t := NewTransaction(sender, recipient, value)

	if sender == MINING_SENDER {
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	}

	if bc.VerifyTransactionSignature(senderPublicKey, s, t) {
		//if bc.CalculateTotalAmount(sender) < value {
		//	log.Println("ERROR: Not enough balance in a wallet")
		//	return false
		//}
		bc.transactionPool = append(bc.transactionPool, t)
		return true
	} else {
		log.Println("ERROR: Verify Transaction")
	}
	return false

}

func (bc *Blockchain) VerifyTransactionSignature(senderPublicKey *ecdsa.PublicKey, s *utils.Signature, t *Transaction) bool {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	return ecdsa.Verify(senderPublicKey, h[:], s.R, s.S)
}

func (bc *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, t := range bc.transactionPool {
		transactions = append(transactions,
			NewTransaction(t.senderBlockchainAddress, t.recipientBlockchainAddress, t.value))
	}
	return transactions
}

func (bc *Blockchain) validProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	if len(ch) == 0 {
		zeros := strings.Repeat("0", difficulty)
		guessBlock := Block{0, nonce, previousHash, transactions}
		guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
		if guessHashStr[:difficulty] == zeros {
			ch <- nonce
		}
	}
	//zeros := strings.Repeat("0", difficulty)
	//guessBlock := Block{0, nonce, previousHash, transactions}
	//guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	//return guessHashStr[:difficulty] == zeros
}

func (bc *Blockchain) ProofOfWork() int {
	transactions := bc.CopyTransactionPool()
	previousHash := bc.LastBlock().Hash()
	nonce := 0
	ch := make(chan int, 1)
	var wg sync.WaitGroup
Loop:
	for {
		select {
		case nonce = <-ch:
			wg.Wait()
			close(ch)
			break Loop
		default:
			wg.Add(1)
			go bc.validProof(nonce, previousHash, transactions, MINING_DIFFICULTY, ch, &wg)
			nonce++
		}
	}
	//for !bc.validProof(nonce, previousHash, transactions, MINING_DIFFICULTY) {
	//	nonce++
	//}
	return nonce
}

func (bc *Blockchain) Mining() bool {
	log.Println("action=mining, status=in_progress")
	bc.AddTransaction(MINING_SENDER, bc.blockchainAddress, MINING_REWARD, nil, nil)
	nonce := bc.ProofOfWork()
	previousHash := bc.LastBlock().Hash()
	bc.CreateBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

func (bc *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, b := range bc.chain {
		for _, t := range b.transactions {
			value := t.value
			if blockchainAddress == t.recipientBlockchainAddress {
				totalAmount += value
			}
			if blockchainAddress == t.senderBlockchainAddress {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}
