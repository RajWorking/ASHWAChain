package main

import (
	"github.com/davecgh/go-spew/spew"
	"log"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

func start() {
	t := time.Now()
	genesisBlock := powBlock{Index: 0, Timestamp: t.String(), PowID: "Genesis ID", Difficulty: PowDifficulty}
	genesisBlock.Hash = powBlockHash(genesisBlock)
	spew.Dump(genesisBlock)
	mutex.Lock()
	powChain = append(powChain, genesisBlock)
	mutex.Unlock()

	comGenesisBlock := comBlock{Index: 0, Timestamp: t.String(), Committee: []string{"Genesis Committee"}}
	comGenesisBlock.Hash = comBlockHash(comGenesisBlock)
	comChain = append(comChain, comGenesisBlock)

	txGenesisBlock := txBlock{Index: 0, Timestamp: t.String(), Transactions: []Tx{},
		ComSign: "Genesis Committee Signature"}
	txGenesisBlock.Hash = txBlockHash(txGenesisBlock)
	txChain = append(txChain, txGenesisBlock)

	log.Fatal(run())
}

func main() {
	start()
}
