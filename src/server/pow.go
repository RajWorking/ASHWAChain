package main

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"net/http"
	"strconv"
	"time"
)

// powBlock represents each block in the powChain.
// The PowID is the one controlled by the node proposing the block
type powBlock struct {
	Index      int
	Timestamp  string
	PowID      string
	PrevHash   string
	Difficulty int
	Nonce      string
	Hash       string
}

// series of validated Blocks
var powChain []powBlock

// PowMessage takes incoming JSON payload for writing some fixed derivative of public key of identity
type PowMessage struct {
	PowID string
}

// write powChain when we receive an http request
func handleGetPowChain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(powChain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

// SHA256 hasing
func powBlockHash(block powBlock) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.PowID + block.PrevHash + block.Nonce
	return calculateHash(record)
}

// make sure PowBlock is valid by checking index, and comparing the hash of the previous PowBlock
func isPowBlockValid(newBlock, oldBlock powBlock) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if powBlockHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

// takes JSON payload as an input for pow id (PowID)
func handleWritePowBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m PowMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newPowBlock := generatePowBlock(powChain[len(powChain)-1], m.PowID)

	//ensure atomicity when adding new block
	mutex.Lock()
	if isPowBlockValid(newPowBlock, powChain[len(powChain)-1]) {
		powChain = append(powChain, newPowBlock)
		spew.Dump(powChain[len(powChain)-1])

		// After every committeeSize blocks, create a new committee and
		// add a block to comChain with new commmittee
		if len(powChain) > CommitteeSize && len(powChain)%Epoch == 1 {
			generateComBlock(len(powChain))
		}
		respondWithJSON(w, r, http.StatusCreated, newPowBlock)
	} else {
		respondWithJSON(w, r, http.StatusForbidden, newPowBlock)
	}
	mutex.Unlock()

}

// create a new powBlock using previous block's hash
func generatePowBlock(oldBlock powBlock, PowID string) powBlock {
	var newBlock powBlock

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PowID = PowID
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = PowDifficulty

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(powBlockHash(newBlock), newBlock.Difficulty) {
			//fmt.Println(powBlockHash(newBlock), " do more work!")
			continue
		} else {
			//fmt.Println(powBlockHash(newBlock), " work done!")
			newBlock.Hash = powBlockHash(newBlock)
			break
		}
	}
	return newBlock
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
