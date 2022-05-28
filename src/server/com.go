package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/* The "Committee" stores a concatenation of IDs from the latest "committeeSize" blocks. */

type comBlock struct {
	Index     int
	Timestamp string
	Committee []string
	Hash      string
	PrevHash  string
}

// series of validated Blocks
var comChain []comBlock

// write comChain when we receive an http request
func handleGetComChain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(comChain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func generateComBlock(height int) {
	var committee []string
	for i := height - CommitteeSize; i < height; i++ {
		committee = append(committee, powChain[i].PowID)
	}
	oldBlock := comChain[len(comChain)-1]
	block := comBlock{len(comChain), time.Now().String(), committee, "", oldBlock.Hash}
	block.Hash = comBlockHash(block)
	comChain = append(comChain, block)
	ioutil.WriteFile("../committee_nodes.txt", []byte(strings.Join(block.Committee, "\n")), 0644)
}

func comBlockHash(block comBlock) string {
	committee := ""
	for i := 0; i < len(block.Committee); i++ {
		committee += block.Committee[i]
	}
	record := strconv.Itoa(block.Index) + block.Timestamp + committee + block.PrevHash
	return calculateHash(record)
}
