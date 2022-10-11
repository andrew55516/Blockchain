package main

import (
	"Blockchain/internal/blockchain"
	"Blockchain/internal/wallet"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("blockchain: ")

}

func main() {
	wM := wallet.NewWallet()
	wA := wallet.NewWallet()
	wB := wallet.NewWallet()

	t := wallet.NewTransaction(wA.PrivateKey(), wA.PublicKey(), wA.BlockchainAddress(), wB.BlockchainAddress(), 1.0)

	bc := blockchain.NewBlockchain(wM.BlockchainAddress())
	isAdded := bc.AddTransaction(wA.BlockchainAddress(), wB.BlockchainAddress(), 2.0, wA.PublicKey(), t.GenerateSignature())
	fmt.Println("Added? ", isAdded)

	bc.Mining()
	bc.Print()

	fmt.Printf("A %f\n", bc.CalculateTotalAmount(wA.BlockchainAddress()))
	fmt.Printf("B %f\n", bc.CalculateTotalAmount(wB.BlockchainAddress()))
	fmt.Printf("A %f\n", bc.CalculateTotalAmount(wM.BlockchainAddress()))
}
