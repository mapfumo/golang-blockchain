package blockchain

import (
	"encoding/hex"
	"fmt"
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
	if DBexists() == false {
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

// AddBlock creates a new block with the given data and adds it to the blockchain.
func (bc *BlockChain) AddBlock(transactions []*Transaction) {
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
	newBlock := CreateBlock(transactions, lastHash)

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

// find all thhe unspent transactions associated with the address
func (bc *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := bc.Iterator()

	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOutdx := range spentTXOs[txID] {
						if spentOutdx == outdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)

				}
			}
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}

		}

		if len(block.PrevHash) == 0 { // reached Genesis block
			break
		}
	}

	return unspentTxs
}

func (bc *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs

}

func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
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
