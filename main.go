package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/mapfumo/golang-blockchain/blockchain"
)

// CommandLine struct to interact with the blockchain
type CommandLine struct {
	bc *blockchain.BlockChain
}

// printUsage prints the usage instructions for the CLI.
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - Add a new block to the blockchain")
	fmt.Println(" print - Prints the blocks on the blockchain")
}

// validateArgs checks if the necessary command-line arguments are provided.
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		// Exit the goroutine to prevent BadgerDB from being improperly garbage collected.
		runtime.Goexit()
	}
}

// addBlock adds a new block with the given data to the blockchain.
func (cli *CommandLine) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Block added")
}

// printChain prints all the blocks in the blockchain.
func (cli *CommandLine) printChain() {
	iter := cli.bc.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("PrevHash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
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

// run parses the command-line arguments and executes the corresponding commands.
func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			cli.printUsage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {
	// Ensure the program exits cleanly to avoid database corruption.
	defer os.Exit(0)

	// Initialize the blockchain.
	bc := blockchain.InitBlockChain()
	// Close the database when the program exits.
	defer bc.Database.Close()

	// Create a CommandLine instance and run the CLI.
	cli := &CommandLine{bc: bc}
	cli.run()
}
