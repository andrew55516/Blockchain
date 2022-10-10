package main

import (
	"Blockchain/internal/blockchain"
	"log"
)

func init() {
	log.SetPrefix("blockchain: ")

}

func main() {
	blockchain := blockchain.NewBlockchain()
	blockchain.Print()
	blockchain.AddTransaction("A", "B", 1.0)
	nonce := blockchain.ProofOfWork()
	blockchain.CreateBlock(nonce, blockchain.LastBlock().Hash())
	blockchain.AddTransaction("C", "D", 3.0006)
	blockchain.AddTransaction("A", "B", 1.0)
	nonce = blockchain.ProofOfWork()
	blockchain.CreateBlock(nonce, blockchain.LastBlock().Hash())
	blockchain.Print()

}
