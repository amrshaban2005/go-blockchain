package main

import (
	"fmt"
	"log"

	"github.com/amrshaban2005/go-blockchain/blockchain"
)

func init() {
	log.SetPrefix("Blockchain Node: ")
}

func main() {
	blockchainAddress := "miners_blockchain_address"
	blockchain := blockchain.NewBlockchain(blockchainAddress)

	blockchain.AddTransactions("Amr", "Moaz", 32)
	blockchain.Mining()

	blockchain.AddTransactions("Ahmed", "AbdelRahman", 30)
	blockchain.AddTransactions("Ali", "Mostafa", 30)
	blockchain.Mining()
	blockchain.Print()

	fmt.Printf("miner %1f\n", blockchain.CalculateTotalAmount(blockchainAddress))
	fmt.Printf("Amr %1f\n", blockchain.CalculateTotalAmount("Amr"))
	fmt.Printf("Moaz %1f\n", blockchain.CalculateTotalAmount("Moaz"))

}
