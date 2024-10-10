package main

import (
	"fmt"
	"log"

	//"github.com/amrshaban2005/go-blockchain/blockchain"
	"github.com/amrshaban2005/go-blockchain/blockchain"
	"github.com/amrshaban2005/go-blockchain/wallet"
)

func init() {
	log.SetPrefix("Blockchain Node: ")
}

func main() {
	Walletminter := wallet.NewWallet()
	walletAlice := wallet.NewWallet()
	walletBob := wallet.NewWallet()
	// fmt.Println(w.PrivateKeyStr())
	// fmt.Println(w.PublicKeyStr())
	// fmt.Println(w.BlockchainAddress())

	t := wallet.NewTransaction(walletAlice.PrivateKey(), walletAlice.PublicKey(), walletAlice.BlockchainAddress(), walletBob.BlockchainAddress(), 32)
	// fmt.Printf("Signature %s", t.GenerateSignature())

	blockchain := blockchain.NewBlockchain(Walletminter.BlockchainAddress())

	isAdded := blockchain.AddTransactions(walletAlice.BlockchainAddress(), walletBob.BlockchainAddress(), 32, walletAlice.PublicKey(), t.GenerateSignature())
	fmt.Println("is transaction verified?", isAdded)
	blockchain.Mining()
	blockchain.Print()

	fmt.Printf("Miner %10f\n", blockchain.CalculateTotalAmount(Walletminter.BlockchainAddress()))
	fmt.Printf("Alice %10f\n", blockchain.CalculateTotalAmount(walletAlice.BlockchainAddress()))
	fmt.Printf("Bob %10f\n", blockchain.CalculateTotalAmount(walletBob.BlockchainAddress()))

	// blockchain.AddTransactions("Ahmed", "AbdelRahman", 30)
	// blockchain.AddTransactions("Ali", "Mostafa", 30)
	// blockchain.Mining()
	// blockchain.Print()

	// fmt.Printf("miner %1f\n", blockchain.CalculateTotalAmount(blockchainAddress))
	// fmt.Printf("Amr %1f\n", blockchain.CalculateTotalAmount("Amr"))
	// fmt.Printf("Moaz %1f\n", blockchain.CalculateTotalAmount("Moaz"))

}
