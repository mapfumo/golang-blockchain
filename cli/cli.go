package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/mapfumo/golang-blockchain/blockchain"
)

// CommandLine struct to interact with the blockchain
type CommandLine struct {}

// printUsage prints the usage instructions for the CLI.
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get the balance of an address")
	fmt.Println(" createblockchain -address ADDRESS - create a new blockchain")
	fmt.Println(" printchain - Prints the blocks on the blockchain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - send amount of tokens from FROM to TO")
}

// validateArgs checks if the necessary command-line arguments are provided.
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		// Exit the goroutine to prevent BadgerDB from being improperly garbage collected.
		runtime.Goexit()
	}
}


// printChain prints all the blocks in the blockchain.
func (cli *CommandLine) printChain() {
	bc := blockchain.ContinueBlockChain("")
	defer bc.Database.Close()
	iter := bc.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		// Break the loop if we reach the genesis block (no previous hash).
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) getBalance(address string) {
	bc := blockchain.ContinueBlockChain(address)
	defer bc.Database.Close()

	balance :=0
	UTXOs := bc.FindUTXO(address)
	for _, out := range UTXOs {
		balance += out.Value
	}
	fmt.Printf("Balance of %s: %d\n", address, balance)
}


func (cli *CommandLine) send(from, to string, amount int) {
	bc := blockchain.ContinueBlockChain(from)
	defer bc.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, bc)
	bc.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Transaction added - Success!")
}

func (cli *CommandLine) createBlockChain(address string) {
	bc := blockchain.InitBlockChain(address)
	bc.Database.Close()
	fmt.Println("Blockchain created - Finshed!")

}

// run parses the command-line arguments and executes the corresponding commands.
func (cli *CommandLine) Run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}