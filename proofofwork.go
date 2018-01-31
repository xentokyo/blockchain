package main

import (
	"math"
	"math/big"
	"bytes"
	"crypto/sha256"
)

const (
	maxNonce = math.MaxInt64
	targetBits = 24
)

type ProofOfWork struct {
	Block *Block
	Target *big.Int
}

func (proofOfWork *ProofOfWork) prepareData(nonce int) []byte {
	return bytes.Join([][]byte{
		proofOfWork.Block.PrevBlockHash,
		proofOfWork.Block.Data,
		IntToHex(proofOfWork.Block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})
}

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

func (proofOfWork *ProofOfWork) Validate() bool {
	var hashInt big.Int

	hash := sha256.Sum256(proofOfWork.prepareData(proofOfWork.Block.Nonce))
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(proofOfWork.Target) == -1
}

func NewProofOfWork(block *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))

	return &ProofOfWork{ block, target }
}
