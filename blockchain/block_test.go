package blockchain

import (
	"bytes"
	"testing"
)

// TestCreateBlock tests the CreateBlock function.
func TestCreateBlock(t *testing.T) {
	data := "test data"
	prevHash := []byte("previous hash")
	block := CreateBlock(data, prevHash)

	if !bytes.Equal(block.Data, []byte(data)) {
		t.Errorf("CreateBlock() Data = %s; want %s", block.Data, data)
	}
	if !bytes.Equal(block.PrevHash, prevHash) {
		t.Errorf("CreateBlock() PrevHash = %x; want %x", block.PrevHash, prevHash)
	}
	if len(block.Hash) == 0 {
		t.Error("CreateBlock() Hash is empty; want non-empty hash")
	}
}

// TestGenesisBlock tests the GenesisBlock function.
func TestGenesisBlock(t *testing.T) {
	block := GenesisBlock()

	if !bytes.Equal(block.Data, []byte("Genesis")) {
		t.Errorf("GenesisBlock() Data = %s; want Genesis", block.Data)
	}
	if len(block.PrevHash) != 0 {
		t.Errorf("GenesisBlock() PrevHash = %x; want empty", block.PrevHash)
	}
	if len(block.Hash) == 0 {
		t.Error("GenesisBlock() Hash is empty; want non-empty hash")
	}
}

// TestInitBlockChain tests the InitBlockChain function.
func TestInitBlockChain(t *testing.T) {
	bc := InitBlockChain()

	if len(bc.Blocks) != 1 {
		t.Errorf("InitBlockChain() Blocks length = %d; want 1", len(bc.Blocks))
	}

	genesisBlock := bc.Blocks[0]
	if !bytes.Equal(genesisBlock.Data, []byte("Genesis")) {
		t.Errorf("InitBlockChain() GenesisBlock Data = %s; want Genesis", genesisBlock.Data)
	}
	if len(genesisBlock.PrevHash) != 0 {
		t.Errorf("InitBlockChain() GenesisBlock PrevHash = %x; want empty", genesisBlock.PrevHash)
	}
	if len(genesisBlock.Hash) == 0 {
		t.Error("InitBlockChain() GenesisBlock Hash is empty; want non-empty hash")
	}
}

// TestAddBlock tests the AddBlock function.
func TestAddBlock(t *testing.T) {
	bc := InitBlockChain()
	bc.AddBlock("Block 1")

	if len(bc.Blocks) != 2 {
		t.Errorf("AddBlock() Blocks length = %d; want 2", len(bc.Blocks))
	}

	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	if !bytes.Equal(lastBlock.Data, []byte("Block 1")) {
		t.Errorf("AddBlock() LastBlock Data = %s; want Block 1", lastBlock.Data)
	}

	prevBlock := bc.Blocks[len(bc.Blocks)-2]
	if !bytes.Equal(lastBlock.PrevHash, prevBlock.Hash) {
		t.Errorf("AddBlock() LastBlock PrevHash = %x; want %x", lastBlock.PrevHash, prevBlock.Hash)
	}
}


