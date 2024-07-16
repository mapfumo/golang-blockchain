package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

/*
Extract the necessary data from the block for hashing.

Initialize a counter (nonce) starting at 0 to vary the input data.

Compute a SHA-256 hash of the block data combined with the current nonce value.

Check if the resulting hash meets the specified difficulty requirements.

Difficulty Requirements:
The hash must have a certain number of leading zeros, which makes finding a valid hash computationally challenging.
This difficulty is dynamically adjusted to regulate the rate of new block creation.
*/

// Difficulty defines the mining difficulty. A higher value makes it more difficult to find a valid hash.
const Difficulty = 12

// ProofOfWork represents a proof-of-work algorithm.
type ProofOfWork struct {
	Block  *Block  // Block is the block being mined.
	Target *big.Int // Target is the hash target that a valid hash must be less than.
}

// NewProofOfWork initializes a new ProofOfWork for a given block.
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	// Shift the target left by (256 - Difficulty) bits to set the difficulty.
	target.Lsh(target, uint(256-Difficulty))
	// Return a new ProofOfWork with the block and target.
	return &ProofOfWork{Block: block, Target: target}
}

// InitData prepares the data for hashing by combining the block's attributes and the nonce.
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,        // Previous block's hash.
			pow.Block.Data,            // Current block's data.
			ToHex(int64(nonce)),       // Nonce converted to a byte slice.
			ToHex(int64(Difficulty)),  // Difficulty level converted to a byte slice.
		},
		[]byte{}) // Separator (none needed here).
	return data
}

// Run is the main method of the ProofOfWork. 
// It loops until a valid hash is found or the nonce reaches the maximum value
// meeting the difficulty target.
func (pow *ProofOfWork) Run() (int, []byte) {
	// intHash is used to hold the hash as a big.Int to facilitate comparison with the target.
	var intHash big.Int
	// hash stores the result of the SHA-256 hash operation.
	var hash [32]byte

	// nonce starts at 0 and will be incremented until a valid hash is found or the maximum value is reached.
	nonce := 0

	// Loop until a valid hash is found or nonce reaches the maximum possible value (math.MaxInt64).
	for nonce < math.MaxInt64 {
		// Prepare the data to be hashed by combining the block's attributes and the current nonce.
		data := pow.InitData(nonce)
		// Compute the SHA-256 hash of the data.
		hash = sha256.Sum256(data)

		// Print the current hash to the console. \r is used to overwrite the same line.
		fmt.Printf("\rHash: %x", hash)

		// Convert the hash to a big.Int to compare it with the target value.
		intHash.SetBytes(hash[:])

		// Compare the hash with the target. If the hash is less than the target, we've found a valid nonce.
		if intHash.Cmp(pow.Target) == -1 {
			break // Exit the loop if a valid nonce is found.
		} else {
			// If the hash is not less than the target, increment the nonce and try again.
			nonce++
		}
	}

	// Print a newline to the console to move to the next line after the final hash output.
	fmt.Println()

	// Return the valid nonce and the corresponding hash as a byte slice.
	return nonce, hash[:]
}


func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(pow.Target) == -1
}

// ToHex converts an int64 number to a byte slice in Big Endian format.
func ToHex(num int64) []byte {
	// Create a new Buffer. Buffers are variable-sized buffers of bytes
	// with Read and Write methods.
	buff := new(bytes.Buffer)

	// Write the binary representation of the number 'num' into the buffer.
	// binary.Write writes the data into the buffer in Big Endian format.
	// Big Endian format means the most significant byte is stored first.
	err := binary.Write(buff, binary.BigEndian, num)

	// Check if there was an error while writing to the buffer.
	// If there was an error, log it and stop the program using log.Panic.
	if err != nil {
		log.Panic(err)
	}

	// Return the byte slice containing the binary representation of 'num'.
	return buff.Bytes()
}
