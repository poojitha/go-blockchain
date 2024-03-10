package main

import (
	"fmt"
	"go-blockchain/blockchain/block"
	"go-blockchain/blockchain/wallet"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {

	wmM := wallet.NewWallet()
	wmA := wallet.NewWallet()
	wmB := wallet.NewWallet()
	t := wallet.NewTransaction(wmA.PrivateKey(), wmA.PublicKey(), wmA.BlockchainAddress(), wmB.BlockchainAddress(), 1.0)

	blockChain := block.NewBlockchain(wmM.BlockchainAddress())
	isAdded := blockChain.AddTransaction(wmA.BlockchainAddress(), wmB.BlockchainAddress(), 1.0, wmA.PublicKey(), t.GenerateSignature())
	fmt.Println("Added?", isAdded)

	blockChain.Print()

	blockChain.Mining()

	blockChain.Print()

	fmt.Printf("A %.1f\n", blockChain.CalculateTotalAmount(wmA.BlockchainAddress()))
	fmt.Printf("B %.1f\n", blockChain.CalculateTotalAmount(wmB.BlockchainAddress()))
	fmt.Printf("M %.1f\n", blockChain.CalculateTotalAmount(wmM.BlockchainAddress()))

}
