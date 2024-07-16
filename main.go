package main

import (
	"fmt"

	"github.com/mapfumo/golang-blockchain/blockchain"
)



func main() {
	bc := blockchain.InitBlockChain()
	bc.AddBlock("Block 1")
	bc.AddBlock("Block 2")
	bc.AddBlock("Block 3")

	for _, block := range bc.Blocks {
		fmt.Printf("PrevHash: %x\t\n", block.PrevHash)
		fmt.Printf("Data: %s\t\t\n", block.Data)
		fmt.Printf("Hash: %x\t\t\n\n", block.Hash)
	}

}
