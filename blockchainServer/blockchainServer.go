package main

import (
	"Blockchain/internal/blockchain"
	"Blockchain/internal/wallet"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
)

var cache map[string]*blockchain.Blockchain = make(map[string]*blockchain.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (bcs *BlockchainServer) Port() uint16 {
	return bcs.port
}

func (bcs *BlockchainServer) GetChain(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bc := bcs.GetBlockchain()
		m, _ := json.Marshal(bc)
		io.WriteString(w, string(m[:]))
	default:
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (bcs *BlockchainServer) Start() {
	http.HandleFunc("/", bcs.GetChain)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(bcs.Port())), nil))
}

func (bcs *BlockchainServer) GetBlockchain() *blockchain.Blockchain {
	bc, ok := cache["blockchain"]
	if !ok {
		minersWallet := wallet.NewWallet()
		bc = blockchain.NewBlockchain(minersWallet.BlockchainAddress(), bcs.port)
		cache["blockchain"] = bc
		log.Printf("private_key: %v\n", minersWallet.PrivateKeyStr())
		log.Printf("publick_key: %v\n", minersWallet.PublicKeyStr())
		log.Printf("blockchain_address: %v\n", minersWallet.BlockchainAddress())
	}
	return bc
}
