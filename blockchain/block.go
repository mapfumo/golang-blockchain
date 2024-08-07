package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash         []byte         // hash of this block - derived from the Data and prevHash
	Transactions []*Transaction // each block needs at least one transaction
	PrevHash     []byte         // hash of the previous block
	Nonce        int

	// Timestamp int
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func CreateBlock(txs []*Transaction, prevHash []byte) *Block {
	b := &Block{[]byte{}, txs, prevHash, 0}
	pow := NewProofOfWork(b)
	nonce, hash := pow.Run()
	b.Nonce = nonce
	b.Hash = hash[:]
	return b
}

func GenesisBlock(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	Handle(err)
	return res.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	Handle(err)
	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
