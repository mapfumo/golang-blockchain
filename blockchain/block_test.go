package blockchain

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

// TestGenesisBlock tests the initialization of the Genesis block.
func TestGenesisBlock(t *testing.T) {
	// Initialize the blockchain
	bc := InitBlockChain()
	// Get the Genesis block
	genesisBlock := bc.Blocks[0]

	// Check if the data of the Genesis block is "Genesis"
	if !bytes.Equal(genesisBlock.Data, []byte("Genesis")) {
		t.Errorf("Expected Genesis block data to be 'Genesis', got '%s'", genesisBlock.Data)
	}

	// Check if the PrevHash of the Genesis block is empty
	if len(genesisBlock.PrevHash) != 0 {
		t.Errorf("Expected Genesis block PrevHash to be empty, got '%x'", genesisBlock.PrevHash)
	}

	// Check if the hash of the Genesis block is correct
	expectedHash := sha256.Sum256([]byte("Genesis"))
	if !bytes.Equal(genesisBlock.Hash, expectedHash[:]) {
		t.Errorf("Expected Genesis block hash to be '%x', got '%x'", expectedHash, genesisBlock.Hash)
	}
}

// TestAddBlock tests adding the first block to the blockchain.
func TestAddBlock(t *testing.T) {
	// Initialize the blockchain
	bc := InitBlockChain()
	// Add a new block with data "First Block"
	bc.AddBlock("First Block")

	// Check if the blockchain length is 2 (Genesis block + first block)
	if len(bc.Blocks) != 2 {
		t.Fatalf("Expected blockchain length to be 2, got %d", len(bc.Blocks))
	}

	// Get the Genesis block and the first block
	genesisBlock := bc.Blocks[0]
	firstBlock := bc.Blocks[1]

	// Check if the data of the first block is "First Block"
	if !bytes.Equal(firstBlock.Data, []byte("First Block")) {
		t.Errorf("Expected first block data to be 'First Block', got '%s'", firstBlock.Data)
	}

	// Check if the PrevHash of the first block matches the hash of the Genesis block
	if !bytes.Equal(firstBlock.PrevHash, genesisBlock.Hash) {
		t.Errorf("Expected first block PrevHash to be '%x', got '%x'", genesisBlock.Hash, firstBlock.PrevHash)
	}

	// Check if the hash of the first block is correct
	expectedHash := sha256.Sum256(append(firstBlock.Data, firstBlock.PrevHash...))
	if !bytes.Equal(firstBlock.Hash, expectedHash[:]) {
		t.Errorf("Expected first block hash to be '%x', got '%x'", expectedHash, firstBlock.Hash)
	}
}

// TestMultipleBlocks tests adding multiple blocks to the blockchain.
func TestMultipleBlocks(t *testing.T) {
	// Initialize the blockchain
	bc := InitBlockChain()
	// Add a new block with data "First Block"
	bc.AddBlock("First Block")
	// Add a new block with data "Second Block"
	bc.AddBlock("Second Block")

	// Check if the blockchain length is 3 (Genesis block + first block + second block)
	if len(bc.Blocks) != 3 {
		t.Fatalf("Expected blockchain length to be 3, got %d", len(bc.Blocks))
	}

	// Get the Genesis block, the first block, and the second block
	firstBlock := bc.Blocks[1]
	secondBlock := bc.Blocks[2]


	// Check if the data of the second block is "Second Block"
	if !bytes.Equal(secondBlock.Data, []byte("Second Block")) {
		t.Errorf("Expected second block data to be 'Second Block', got '%s'", secondBlock.Data)
	}

	// Check if the PrevHash of the second block matches the hash of the first block
	if !bytes.Equal(secondBlock.PrevHash, firstBlock.Hash) {
		t.Errorf("Expected second block PrevHash to be '%x', got '%x'", firstBlock.Hash, secondBlock.PrevHash)
	}

	// Check if the hash of the second block is correct
	expectedHash := sha256.Sum256(append(secondBlock.Data, secondBlock.PrevHash...))
	if !bytes.Equal(secondBlock.Hash, expectedHash[:]) {
		t.Errorf("Expected second block hash to be '%x', got '%x'", expectedHash, secondBlock.Hash)
	}
}
