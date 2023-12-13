package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const (
	DB_FILE                = "blockchain.db"
	BLOCKS_BUCKET          = "blocks"
	GENESIS_COIN_BASE_DATA = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// AddBlock saves provided data as a block in the blockchain
func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		fmt.Println("Error " + err.Error())
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			fmt.Println("Error " + err.Error())
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			fmt.Println("Error " + err.Error())
		}

		bc.tip = newBlock.Hash

		return nil
	})
}

func dbExists() bool {
	if _, err := os.Stat(DB_FILE); os.IsNotExist(err) {
		return false
	}

	return true
}

func NewBlockchain(address string) *Blockchain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(DB_FILE, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		tip = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}

// CreateBlockchain creates a new blockchain DB
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(DB_FILE, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, GENESIS_COIN_BASE_DATA)
		genesis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(BLOCKS_BUCKET))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}
