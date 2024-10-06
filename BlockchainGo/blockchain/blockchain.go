package blockchain

import (
    "bytes"
    "encoding/gob"
	"github.com/dgraph-io/badger"
)

const(
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB

}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	Handle(err)

	err := db.Update(func(txn *badger.Txn) error{
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound{
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Gebesus proved")
			err = txn.Set(genesis.Hash, genesis.Serialize())
			err := txn.Get([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.Value()
			return err
		}
	})

	Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

