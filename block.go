package main

import (
	"time"
)

type Block struct {
	Timestamp int64
	Data []byte
	PrevBlockHash []byte
	Hash []byte
	Nonce int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{ time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0 }
	proofOfWork := NewProofOfWork(block)
	nonce, hash := proofOfWork.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("This is the genesis block", []byte{})
}
