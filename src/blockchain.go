package main

import (
	"fmt"

	"github.com/boltdb/bolt"
)

const (
	DB_FILE       = "../../blockchain.db"
	BLOCKS_BUCKET = "blocks"
)

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// AddBlock saves provided data as a block in the blockchain
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		fmt.Println("Error " + err.Error())
	}

	newBlock := NewBlock(data, lastHash)

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

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(DB_FILE, 0600, nil)
	if err != nil {
		fmt.Println("Error " + err.Error())
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(BLOCKS_BUCKET))
			if err != nil {
				fmt.Println("Error " + err.Error())
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				fmt.Println("Error " + err.Error())
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				fmt.Println("Error " + err.Error())
			}

			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})
	bc := Blockchain{tip, db}

	return &bc
}
