package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dangitsdavid/go-blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	for {
		// Add a new block with the current timestamp as the data
		chain.AddBlock(fmt.Sprintf("Block added at %v", time.Now().Format(time.RFC3339)))

		// Sleep for 2 seconds before adding the next block
		time.Sleep(2 * time.Second)

		// print out block data
		for _, block := range chain.Blocks {
			fmt.Printf("Previous Hash: %x\n", block.PrevHash)
			fmt.Printf("Data in Block: %s\n", block.Data)
			fmt.Printf("Hash: %x\n", block.Hash)

			pow := blockchain.NewProof(block)
			fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
			fmt.Println()
		}
	}
}
