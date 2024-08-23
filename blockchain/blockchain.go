package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"

	badger "github.com/dgraph-io/badger"
)

// dbPath is the directory where the BadgerDB database files will be stored.
const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from Genesis Block"
)

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

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}



// InitBlockChain initializes a new blockchain or loads an existing one from the database.
func InitBlockChain(address string) *BlockChain {
	fmt.Println("Initializing blockchain...")
	var lastHash []byte

	if DBexists() {
		fmt.Println("Blockchain already exists - DB exists")
		runtime.Goexit()
	}

	// Set up BadgerDB options with the database path.
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil
	opts.ValueDir = dbPath

	// Open the BadgerDB database.
	db, err := badger.Open(opts)
	Handle(err)

	// Update the database to initialize or load the blockchain.
	err = db.Update(func(txn *badger.Txn) error {
		fmt.Println("Creating genesis block...")
		cbtx := CoinBaseTx(address, genesisData) // the address that will be rewarded 100 tokens
		genesis := GenesisBlock(cbtx)
		fmt.Println("Genesis block created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)
		// Set the lastHash key to point to the genesis block.
		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash
		return err
	})

	Handle(err)

	// Return a new BlockChain instance with the last hash and database.
	return &BlockChain{
		LastHash: lastHash,
		Database: db,
	}
}

func ContinueBlockChain(address string) *BlockChain {
	if !DBexists() {
		fmt.Println("Blockchain does not exist - Create one!t")
		runtime.Goexit()
	}
	var lastHash []byte

	// Set up BadgerDB options with the database path.
	opts := badger.DefaultOptions(dbPath)
	opts.ValueDir = dbPath

	// Open the BadgerDB database.
	db, err := badger.Open(opts)
	Handle(err)

	// Update the database to initialize or load the blockchain.
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)
		return err
	})
	Handle(err)

	bc := &BlockChain{
		LastHash: lastHash,
		Database: db,
	}

	return bc
}



func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	for _, tx := range transactions {
		if !chain.VerifyTransaction(tx) {
			log.Panic("Invalid Transaction...")
		}
	}

	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})
	Handle(err)

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		encodedBlock, err := item.ValueCopy(nil)
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}

func (chain *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

func (chain *BlockChain) FindUTXO(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (chain *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction does not exist")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.ID)
		Handle(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}