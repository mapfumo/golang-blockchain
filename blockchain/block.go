package blockchain

import (
	"bytes"
	"crypto/sha256"
)

type Block struct {
	Hash      []byte	// hash of this block - derived from the Data and prevHash
	PrevHash  []byte	// hash of the previous block
	Data      []byte
	// Nonce     int
	// Timestamp int
}

type BlockChain struct {
	Blocks []*Block
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	b := &Block{
		Data:      []byte(data),
		PrevHash:  prevHash,
	}
	b.DeriveHash()
	return b
}

func GenesisBlock() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	bc := &BlockChain{
		Blocks: []*Block{GenesisBlock()},
	}
	return bc
}

func (bc *BlockChain) AddBlock(data string) {
	prveBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := CreateBlock(data, prveBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)
}