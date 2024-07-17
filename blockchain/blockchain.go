package blockchain

import (
	"fmt"

	badger "github.com/dgraph-io/badger"
)

// dbPath is the directory where the BadgerDB database files will be stored.
const dbPath = "./tmp/blocks"

// BlockChain represents the blockchain with the last block's hash and the database instance.
type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

// BlockChainIterator allows iterating over the blockchain blocks.
type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// InitBlockChain initializes a new blockchain or loads an existing one from the database.
func InitBlockChain() *BlockChain {
	var lastHash []byte

	// Set up BadgerDB options with the database path.
	opts := badger.DefaultOptions(dbPath)
	opts.ValueDir = dbPath

	// Open the BadgerDB database.
	db, err := badger.Open(opts)
	Handle(err)

	// Update the database to initialize or load the blockchain.
	err = db.Update(func(txn *badger.Txn) error {
		// Check if the lastHash key exists (i.e., if we already have a blockchain).
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			// If the key does not exist, create the genesis block.
			fmt.Println("No existing blockchain found")
			genesis := GenesisBlock()
			fmt.Println("Genesis block created")

			// Serialize the genesis block and store it in the database.
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			// Set the lastHash key to point to the genesis block.
			err = txn.Set([]byte("lh"), genesis.Hash)

			// Update the lastHash variable to the genesis block's hash.
			lastHash = genesis.Hash

			return err
		} else {
			// If the key exists, load the last hash from the database.
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.ValueCopy(nil)
			return err
		}
	})

	Handle(err)

	// Return a new BlockChain instance with the last hash and database.
	return &BlockChain{
		LastHash: lastHash,
		Database: db,
	}
}

// AddBlock creates a new block with the given data and adds it to the blockchain.
func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte

	// View the current last hash from the database.
	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)
		return err
	})

	Handle(err)

	// Create a new block with the current last hash.
	newBlock := CreateBlock(data, lastHash)

	// Update the database with the new block and update the last hash.
	err = bc.Database.Update(func(txn *badger.Txn) error {
		err = txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		// Update the blockchain's last hash to the new block's hash.
		bc.LastHash = newBlock.Hash
		return err
	})

	Handle(err)
}

// Iterator creates a new iterator for the blockchain starting from the last block.
func (bc *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{
		CurrentHash: bc.LastHash,
		Database:    bc.Database,
	}

	return iter
}

// Next retrieves the next block in the blockchain using the iterator.
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	// View the current block using the current hash in the iterator.
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.ValueCopy(nil)
		block = Deserialize(encodedBlock)

		return err
	})

	Handle(err)

	// Update the iterator's current hash to the previous hash of the current block.
	iter.CurrentHash = block.PrevHash

	return block
}
