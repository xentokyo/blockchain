package main

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

const (
	maxNonce   = math.MaxInt64
	targetBits = 24
)

// ProofOfWork struct
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func (proofOfWork *ProofOfWork) prepareData(nonce int) []byte {
	return bytes.Join([][]byte{
		proofOfWork.Block.PrevBlockHash,
		proofOfWork.Block.HashTransactions(),
		IntToHex(proofOfWork.Block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})
}

// Run executes proof of work
func (proofOfWork *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < maxNonce {
		hash = sha256.Sum256(proofOfWork.prepareData(nonce))
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(proofOfWork.Target) == -1 {
			break

		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

// Validate returns validity of the proof of work
func (proofOfWork *ProofOfWork) Validate() bool {
	var hashInt big.Int

	hash := sha256.Sum256(proofOfWork.prepareData(proofOfWork.Block.Nonce))
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(proofOfWork.Target) == -1
}

// NewProofOfWork returns new proof of work
func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return &ProofOfWork{block, target}
}
