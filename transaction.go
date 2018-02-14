package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const (
	subsidy = 10
)

// Transaction struct
type Transaction struct {
	ID      []byte
	Outputs []TxOutput
	Inputs  []TxInput
}

// IsCoinbase checks whether the transaction is coinbase
func (transaction Transaction) IsCoinbase() bool {
	return len(transaction.Inputs) == 1 && len(transaction.Inputs[0].TransactionID) == 0 && transaction.Inputs[0].OutputIndex == -1
}

func (transaction Transaction) setID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(transaction)

	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	transaction.ID = hash[:]
}

// TxOutput struct
type TxOutput struct {
	Value        int
	ScriptPubKey string
}

// TxInput struct
type TxInput struct {
	TransactionID []byte
	OutputIndex   int
	ScriptSig     string
}

// CanUnlockOutputWith checks whether the address initiated the transaction
func (input *TxInput) CanUnlockOutputWith(unlockingData string) bool {
	return input.ScriptSig == unlockingData
}

// CanBeUnlockedWith checks if the output can be unlocked with the provided data
func (output *TxOutput) CanBeUnlockedWith(unlockingData string) bool {
	return output.ScriptPubKey == unlockingData
}

// NewCoinbaseTx creates a new coinbase transaction
func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to %s", to)
	}

	input := TxInput{[]byte{}, -1, data}
	output := TxOutput{subsidy, to}
	transaction := Transaction{nil, []TxOutput{output}, []TxInput{input}}

	return &transaction
}

// NewUnspentTransaction creates a transaction
func NewUnspentTransaction(from, to string, amount int, blockchain *Blockchain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := blockchain.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("Not enough funds")
	}

	for encodedID, outputs := range validOutputs {
		transactionID, err := hex.DecodeString(encodedID)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outputs {
			input := TxInput{transactionID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	transaction := Transaction{nil, outputs, inputs}
	transaction.setID()

	return &transaction
}
