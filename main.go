package main

import (
	"fmt"
)

func main() {
	blockchain := NewBlockchain()
	blockchain.AddBlock("Send BTC to Xen")
	blockchain.AddBlock("Send BTC to Mink")

	for _, block := range blockchain.blocks {
		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
