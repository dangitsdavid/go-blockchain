package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/dangitsdavid/go-blockchain/blockchain"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

// usage instructions for cli
func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - Add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

// cli args validation, requires 2 args: add and print
func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		// initiate a shutdown to prevent corrupted database
		runtime.Goexit()
	}
}

// add block method in cli
func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

// print method in cli
func (cli *CommandLine) printChain() {
	iter := cli.blockchain.Iterator()

	// print all previous block data in blockchain
	for {
		block := iter.Next()

		fmt.Printf("Prev. Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		// break out of loop if no more hashes
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

// run which calls all cli methods
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

	// check if addBlockCmd is parsed successfully
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	// check if printChainCmd is parsed successfully
	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {
	// failsafe exit despite whether go main function finishes run
	defer os.Exit(0)

	chain := blockchain.InitBlockChain()

	// properly close database before main function ends
	defer chain.Database.Close()

	cli := CommandLine{chain}
	cli.run()
}
