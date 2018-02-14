package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const (
	dbFile       = "blockchain.db"
	blocksBucket = "blocks"
	genesisData  = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

// Blockchain implements interactions with a DB
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// MineBlock mineds a new block with the provided transactions
func (blockchain *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := blockchain.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		lastHash = bucket.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = blockchain.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		err := bucket.Put(newBlock.Hash, newBlock.Serialize())

		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		blockchain.tip = newBlock.Hash

		return nil
	})
}

// FindUnspentTransactions returns a list of transactions containing unspent outputs
func (blockchain *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTransactions []Transaction
	spentOutputIndexes := make(map[string][]int)
	iterator := blockchain.Iterator()

	for {
		block := iterator.Next()

		for _, transaction := range block.Transactions {
			transactionID := hex.EncodeToString(transaction.ID)

		Outputs:
			for index, output := range transaction.Outputs {
				if spentOutputIndexes[transactionID] != nil {
					for _, spentOutputIndex := range spentOutputIndexes[transactionID] {
						if spentOutputIndex == index {
							continue Outputs
						}
					}
				}

				if output.CanBeUnlockedWith(address) {
					unspentTransactions = append(unspentTransactions, *transaction)
				}
			}

			if !transaction.IsCoinbase() {
				for _, input := range transaction.Inputs {
					if input.CanUnlockOutputWith(address) {
						inputTransactionID := hex.EncodeToString(input.TransactionID)
						spentOutputIndexes[inputTransactionID] = append(spentOutputIndexes[inputTransactionID], input.OutputIndex)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTransactions
}

// FindUnspentTransactionOutputs finds & returns all unspent transaction outouts
func (blockchain *Blockchain) FindUnspentTransactionOutputs(address string) []TxOutput {
	var unspentOutputs []TxOutput
	unspentTransactions := blockchain.FindUnspentTransactions(address)

	for _, transaction := range unspentTransactions {
		for _, output := range transaction.Outputs {
			if output.CanBeUnlockedWith(address) {
				unspentOutputs = append(unspentOutputs, output)
			}
		}
	}

	return unspentOutputs
}

// FindSpendableOutputs finds & returns unspent outputs to reference in inputs
func (blockchain *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTransactions := blockchain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, transaction := range unspentTransactions {
		transactionID := hex.EncodeToString(transaction.ID)

		for index, output := range transaction.Outputs {
			if output.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += output.Value
				unspentOutputs[transactionID] = append(unspentOutputs[transactionID], index)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// Iterator returns struct specified the blockchain
func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.tip, blockchain.db}
}

// Next returns next block in the blockchain
func (iterator *BlockchainIterator) Next() *Block {
	var block *Block

	err := iterator.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(iterator.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	iterator.currentHash = block.PrevBlockHash

	return block
}

// NewBlockchain returns new struct of blockchain
func NewBlockchain() *Blockchain {
	if !dbExists() {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		tip = bucket.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)

	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		transaction := NewCoinbaseTx(address, genesisData)
		genesisBlock := NewGenesisBlock(transaction)

		bucket, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = bucket.Put([]byte("l"), genesisBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		tip = genesisBlock.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}
