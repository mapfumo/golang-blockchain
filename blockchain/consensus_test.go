package blockchain

import (
	"bytes"
	"math/big"
	"testing"
)

// TestToHex tests the ToHex function.
func TestToHex(t *testing.T) {
	num := int64(12345)
	expected := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x30, 0x39}

	result := ToHex(num)
	if !bytes.Equal(result, expected) {
		t.Errorf("ToHex(%d) = %x; want %x", num, result, expected)
	}
}
// TestNewProofOfWork tests the NewProofOfWork function.
func TestNewProofOfWork(t *testing.T) {
	block := &Block{PrevHash: []byte{}, Data: []byte("test block")}
	pow := NewProofOfWork(block)

	expectedTarget := big.NewInt(1)
	expectedTarget.Lsh(expectedTarget, uint(256-Difficulty))

	if pow.Block != block {
		t.Errorf("NewProofOfWork() Block = %v; want %v", pow.Block, block)
	}
	if pow.Target.Cmp(expectedTarget) != 0 {
		t.Errorf("NewProofOfWork() Target = %v; want %v", pow.Target, expectedTarget)
	}
}

// TestInitData tests the InitData method.
func TestInitData(t *testing.T) {
	block := &Block{PrevHash: []byte("previous hash"), Data: []byte("test block")}
	pow := NewProofOfWork(block)

	nonce := 0
	data := pow.InitData(nonce)

	expected := bytes.Join(
		[][]byte{
			block.PrevHash,
			block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	if !bytes.Equal(data, expected) {
		t.Errorf("InitData(%d) = %x; want %x", nonce, data, expected)
	}
}

// TestRun tests the Run method.
func TestRun(t *testing.T) {
	block := &Block{PrevHash: []byte("previous hash"), Data: []byte("test block")}
	pow := NewProofOfWork(block)

	nonce, hash := pow.Run()

	var intHash big.Int
	intHash.SetBytes(hash)

	if intHash.Cmp(pow.Target) != -1 {
		t.Errorf("Run() hash = %x; does not meet the target %x", hash, pow.Target)
	}

	if nonce < 0 {
		t.Errorf("Run() nonce = %d; invalid nonce", nonce)
	}
}

// TestValidate tests the Validate method.
func TestValidate(t *testing.T) {
	block := &Block{PrevHash: []byte("previous hash"), Data: []byte("test block")}
	pow := NewProofOfWork(block)

	nonce, _ := pow.Run()
	block.Nonce = nonce

	if !pow.Validate() {
		t.Errorf("Validate() = false; want true")
	}
}

