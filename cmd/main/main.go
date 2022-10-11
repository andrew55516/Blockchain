package main

import (
	"Blockchain/internal/blockchain"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("blockchain: ")

}

func main() {
	myBlockchainAddress := "my_blockchain_address"
	blockchain := blockchain.NewBlockchain(myBlockchainAddress)
	blockchain.Print()
	blockchain.AddTransaction("A", "B", 1.0)
	blockchain.Mining()
	blockchain.AddTransaction("C", "D", 3.0006)
	blockchain.AddTransaction("A", "B", 1.0)
	blockchain.Mining()
	blockchain.Print()
	fmt.Printf("my %f\n", blockchain.CalculateTotalAmount(myBlockchainAddress))
	fmt.Printf("A %f\n", blockchain.CalculateTotalAmount("A"))
	fmt.Printf("B %f\n", blockchain.CalculateTotalAmount("B"))
	fmt.Printf("C %f\n", blockchain.CalculateTotalAmount("C"))
	fmt.Printf("D %f\n", blockchain.CalculateTotalAmount("D"))
}
