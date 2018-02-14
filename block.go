package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

// Block struct
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// Serialize returns serialized []byte from Block
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// HashTransactions returns a hash of the transactions in the block
func (block *Block) HashTransactions() []byte {
	var hashes [][]byte

	for _, transaction := range block.Transactions {
		hashes = append(hashes, transaction.ID)
	}

	hash := sha256.Sum256(bytes.Join(hashes, []byte{}))

	return hash[:]
}

// NewBlock returns new block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	proofOfWork := NewProofOfWork(block)
	nonce, hash := proofOfWork.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock returns new genesis block
func NewGenesisBlock(transaction *Transaction) *Block {
	return NewBlock([]*Transaction{transaction}, []byte{})
}

// DeserializeBlock returns deserialized Block from []byte
func DeserializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
