package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

// CLI struct
type CLI struct{}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("	checkBalance -address ADDRESS - check balance of address")
	fmt.Println("	createBlockchain -address ADDRESS - create a blockchain & send genegis block reward to address")
	fmt.Println("	send -from FROM -to TO -amount AMOUNT - send amount of coins from address to one")
	fmt.Println("	printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) checkBalance(address string) {
	blockchain := NewBlockchain()
	defer blockchain.db.Close()

	balance := 0
	outputs := blockchain.FindUnspentTransactionOutputs(address)

	for _, output := range outputs {
		balance += output.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CLI) createBlockchain(address string) {
	blockchain := CreateBlockchain(address)
	blockchain.db.Close()
	fmt.Println("Created!")
}

func (cli *CLI) send(from, to string, amount int) {
	blockchain := NewBlockchain()
	defer blockchain.db.Close()

	transaction := NewUnspentTransaction(from, to, amount, blockchain)

	blockchain.MineBlock([]*Transaction{transaction})
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	blockchain := NewBlockchain()
	defer blockchain.db.Close()

	iterator := blockchain.Iterator()

	for {
		block := iterator.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		proofOfWork := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(proofOfWork.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

// Run provides blockchain actions
func (cli *CLI) Run() {
	cli.validateArgs()

	checkBalanceCmd := flag.NewFlagSet("checkBalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	checkBalanceAddress := checkBalanceCmd.String("address", "", "This address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "This address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "checkBalance":
		err := checkBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "createBlockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	default:
		cli.printUsage()
		os.Exit(1)
	}

	if checkBalanceCmd.Parsed() {
		if *checkBalanceAddress == "" {
			checkBalanceCmd.Usage()
			os.Exit(1)
		}

		cli.checkBalance(*checkBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}

		cli.createBlockchain(*createBlockchainAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
