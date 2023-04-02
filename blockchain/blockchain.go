package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// method initializes the genesis block
func InitBlockChain() *BlockChain {
	var lastHash []byte

	// loads badgerdb with default options/optimizations
	opts := badger.DefaultOptions
	// opts.Dir stores keys and metadata
	opts.Dir = dbPath
	// opts.ValueDir stores all the values
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	// initializing genesis block, Update function since db write capabilities needed
	err = db.Update(func(txn *badger.Txn) error {
		// check whether blockchain exists, create genesis and proof
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			// if existing blockchain exists, return lastHash
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.Value()
			return err
		}
	})
	Handle(err)

	// create new blockchain
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

// method to add block to chain
func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	// readonly to check lastHash
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)

		lastHash, err = item.Value()

		return err
	})
	Handle(err)

	// create a new block with the lastHash
	newBlock := CreateBlock(data, lastHash)

	// do a read/write to set the new block's hash
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

// iterator method allows us to view previous blocks
func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

// next function to work with iterator method
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		encodedBlock, err := item.Value()
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
