package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Tx struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value int    `json:"value"`
}

// The "ComSign" is the signature by the Committee.
// Currently, "ComSign' stores concatenation of all the IDs as a string.
type txBlock struct {
	Index        int
	Timestamp    string
	Transactions []Tx
	ComSign      string
	Hash         string
	PrevHash     string
}

// series of validated Blocks
var txChain []txBlock

// list of pending transactions
var pendingTxs []Tx

// recieve transactions, create signed Tx block and append to txChain
func handleWriteTransaction(w http.ResponseWriter, r *http.Request) {
	var m Tx
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	// Check transaction data validity and add the tx to pending txs list
	pendingTxs = append(pendingTxs, m)

	// This happens in PBFT: Create a block of transactions and sign by the current committee
	if len(pendingTxs) >= TxBlockSize {
		generateTxBlock()
	}

	respondWithJSON(w, r, http.StatusCreated, m)
}

func handleGetTxChain(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	//index, err := strconv.Atoi(params["index"])

	bytes, err := json.MarshalIndent(txChain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(bytes))
}

func generateTxBlock() {
	oldBlock := txChain[len(txChain)-1]
	comSign := ""
	committee := comChain[len(comChain)-1].Committee
	for i := 0; i < len(committee); i++ {
		comSign += committee[i]
	}

	block := txBlock{len(txChain), time.Now().String(), pendingTxs, comSign, "", oldBlock.Hash}
	pendingTxs = []Tx{}
	block.Hash = txBlockHash(block)
	txChain = append(txChain, block)
}

func txBlockHash(block txBlock) string {
	txs := ""
	for i := 0; i < len(block.Transactions); i++ {
		txs += fmt.Sprint(block.Transactions[i])
	}
	record := strconv.Itoa(block.Index) + block.Timestamp + txs + block.ComSign + block.PrevHash
	return calculateHash(record)
}
