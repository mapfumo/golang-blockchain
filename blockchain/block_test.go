package blockchain

import (
	"bytes"
	"fmt"
	"testing"
)

func TestCreateBlock(t *testing.T) {
	data := "test data"
	prevHash := []byte("previous hash")
	block := CreateBlock(data, prevHash)

	if !bytes.Equal(block.Data, []byte(data)) {
		t.Errorf("expected block data to be %s, got %s", data, block.Data)
	}
	if !bytes.Equal(block.PrevHash, prevHash) {
		t.Errorf("expected block previous hash to be %x, got %x", prevHash, block.PrevHash)
	}
	if len(block.Hash) == 0 {
		t.Errorf("expected block hash to be set, got empty hash")
	}
	if block.Nonce == 0 {
		t.Errorf("expected block nonce to be set, got 0")
	}
}

func TestGenesisBlock(t *testing.T) {
	block := GenesisBlock()

	if string(block.Data) != "Genesis" {
		t.Errorf("expected Genesis block data to be 'Genesis', got %s", block.Data)
	}
	if len(block.PrevHash) != 0 {
		t.Errorf("expected Genesis block previous hash to be empty, got %x", block.PrevHash)
	}
	if len(block.Hash) == 0 {
		t.Errorf("expected Genesis block hash to be set, got empty hash")
	}
	if block.Nonce == 0 {
		t.Errorf("expected Genesis block nonce to be set, got 0")
	}
}

func TestSerializeDeserialize(t *testing.T) {
	data := "serialize test"
	prevHash := []byte("prev hash")
	block := CreateBlock(data, prevHash)

	serializedBlock := block.Serialize()
	deserializedBlock := Deserialize(serializedBlock)

	if !bytes.Equal(deserializedBlock.Data, block.Data) {
		t.Errorf("expected deserialized block data to be %s, got %s", block.Data, deserializedBlock.Data)
	}
	if !bytes.Equal(deserializedBlock.PrevHash, block.PrevHash) {
		t.Errorf("expected deserialized block previous hash to be %x, got %x", block.PrevHash, deserializedBlock.PrevHash)
	}
	if !bytes.Equal(deserializedBlock.Hash, block.Hash) {
		t.Errorf("expected deserialized block hash to be %x, got %x", block.Hash, deserializedBlock.Hash)
	}
	if deserializedBlock.Nonce != block.Nonce {
		t.Errorf("expected deserialized block nonce to be %d, got %d", block.Nonce, deserializedBlock.Nonce)
	}
}

func TestHandle(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected Handle to panic on error, but it did not")
		}
	}()

	// This should cause Handle to panic
	Handle(fmt.Errorf("test error"))
}
