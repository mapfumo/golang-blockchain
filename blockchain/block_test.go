package blockchain

import (
	"bytes"
	"testing"
)

// Mock Transaction struct for testing purposes
type MockTransaction struct {
	ID []byte
}

func (tx *MockTransaction) Serialize() []byte {
	return tx.ID
}

func TestBlock_SerializeDeserialize(t *testing.T) {
	// Create a sample block with mock transactions
	txs := []*Transaction{
		&Transaction{ID: []byte("tx1")},
		&Transaction{ID: []byte("tx2")},
	}
	block := CreateBlock(txs, []byte("prev-hash"), 1)

	// Serialize the block
	serializedBlock := block.Serialize()

	// Deserialize the block
	deserializedBlock := Deserialize(serializedBlock)

	// Verify that the deserialized block matches the original block
	if deserializedBlock.Timestamp != block.Timestamp {
		t.Errorf("expected Timestamp %d, got %d", block.Timestamp, deserializedBlock.Timestamp)
	}
	if !bytes.Equal(deserializedBlock.Hash, block.Hash) {
		t.Errorf("expected Hash %x, got %x", block.Hash, deserializedBlock.Hash)
	}
	if deserializedBlock.Nonce != block.Nonce {
		t.Errorf("expected Nonce %d, got %d", block.Nonce, deserializedBlock.Nonce)
	}
	if deserializedBlock.Height != block.Height {
		t.Errorf("expected Height %d, got %d", block.Height, deserializedBlock.Height)
	}
	if !bytes.Equal(deserializedBlock.PrevHash, block.PrevHash) {
		t.Errorf("expected PrevHash %x, got %x", block.PrevHash, deserializedBlock.PrevHash)
	}
	if len(deserializedBlock.Transactions) != len(block.Transactions) {
		t.Errorf("expected %d transactions, got %d", len(block.Transactions), len(deserializedBlock.Transactions))
	}
}

func TestCreateBlock(t *testing.T) {
	txs := []*Transaction{
		&Transaction{ID: []byte("tx1")},
		&Transaction{ID: []byte("tx2")},
	}
	prevHash := []byte("prev-hash")
	block := CreateBlock(txs, prevHash, 1)

	if block.Timestamp == 0 {
		t.Errorf("expected valid timestamp, got 0")
	}
	if !bytes.Equal(block.PrevHash, prevHash) {
		t.Errorf("expected PrevHash %x, got %x", prevHash, block.PrevHash)
	}
	if block.Height != 1 {
		t.Errorf("expected Height 1, got %d", block.Height)
	}
	if len(block.Transactions) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(block.Transactions))
	}
}

func TestGenesis(t *testing.T) {
	coinbaseTx := &Transaction{ID: []byte("coinbase-tx")}
	block := Genesis(coinbaseTx)

	if block.Height != 0 {
		t.Errorf("expected Height 0, got %d", block.Height)
	}
	if len(block.Transactions) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(block.Transactions))
	}
	if !bytes.Equal(block.Transactions[0].ID, coinbaseTx.ID) {
		t.Errorf("expected Transaction ID %x, got %x", coinbaseTx.ID, block.Transactions[0].ID)
	}
}

func TestHashTransactions(t *testing.T) {
	txs := []*Transaction{
		&Transaction{ID: []byte("tx1")},
		&Transaction{ID: []byte("tx2")},
	}
	block := CreateBlock(txs, []byte("prev-hash"), 1)
	merkleRoot := block.HashTransactions()

	expectedRoot := NewMerkleTree([][]byte{[]byte("tx1"), []byte("tx2")}).RootNode.Data
	if !bytes.Equal(merkleRoot, expectedRoot) {
		t.Errorf("expected Merkle Root %x, got %x", expectedRoot, merkleRoot)
	}
}

func TestBlock_RunProofOfWork(t *testing.T) {
	txs := []*Transaction{
		&Transaction{ID: []byte("tx1")},
	}
	block := CreateBlock(txs, []byte("prev-hash"), 1)

	// Check if a valid proof of work was created
	pow := NewProofOfWork(block)
	isValid := pow.Validate()

	if !isValid {
		t.Errorf("proof of work validation failed")
	}
}
