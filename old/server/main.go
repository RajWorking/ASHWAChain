package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

const powDifficulty = 1
const committeeSize = 5
const txBlockSize = 3

// powBlock represents each 'item' in the powChain.
// The PowID is the one controlled by the node proposing the block
type powBlock struct {
	Index      int
	Timestamp  string
	PowID      string
	Hash       string
	PrevHash   string
	Difficulty int
	Nonce      string
}

// The "Committee" stores a concatination of IDs from the latest "committeeSize" blocks.
type comBlock struct {
	Index     int
	Timestamp string
	Committee []string
	Hash      string
	PrevHash  string
}

// The "ComSign" is the signature by the Committee.
// Currently "ComSign' stores concatination of all the IDs as a string.
type txBlock struct {
	Index        int
	Timestamp    string
	Transactions []string
	ComSign      string
	Hash         string
	PrevHash     string
}

// <name>Chain is a series of validated <name>Blocks
var powChain []powBlock
var comChain []comBlock
var txChain []txBlock

var pendingTxs []string

// PowMessage takes incoming JSON payload for writing some fixed derivative of public key of identity
type PowMessage struct {
	PowID string
}

var mutex = &sync.Mutex{}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()
		genesisBlock := powBlock{}
		genesisBlock = powBlock{0, t.String(), "Genesis ID", powBlockHash(genesisBlock), "", powDifficulty, ""}
		spew.Dump(genesisBlock)

		mutex.Lock()
		powChain = append(powChain, genesisBlock)
		mutex.Unlock()

		comGenesisBlock := comBlock{}
		comGenesisBlock = comBlock{0, t.String(), []string{"Genesis Committee"}, comBlockHash(comGenesisBlock), ""}
		comChain = append(comChain, comGenesisBlock)

		txGenesisBlock := txBlock{}
		txGenesisBlock = txBlock{0, t.String(), []string{"Genesis Transaction"},
			"Genesis Committee Signature", txBlockHash(txGenesisBlock), ""}
		txChain = append(txChain, txGenesisBlock)
	}()
	log.Fatal(run())

}

// web server
func run() error {
	mux := makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetPowChain).Methods("GET")
	muxRouter.HandleFunc("/tx/{txData}", handleTransactions).Methods("GET")
	muxRouter.HandleFunc("/committee", handleGetComChain).Methods("GET")

	muxRouter.HandleFunc("/", handleWritePowBlock).Methods("POST")
	return muxRouter
}

// recieve transactions, create signed Tx block and append to txChain
func handleTransactions(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	txData := params["txData"]

	// Check transaction data validity and add the tx to pending txs list
	pendingTxs = append(pendingTxs, txData)

	// This happens in PBFT: Create a block of transactions and sign by the current committee
	if len(pendingTxs) >= txBlockSize {
		txBlock := generateTxBlock()
		txChain = append(txChain, txBlock)
	}

	bytes, err := json.MarshalIndent(txChain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(bytes))
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

// write comChain when we receive an http request
func handleGetComChain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(comChain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
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

	//ensure atomicity when creating new block
	mutex.Lock()
	newPowBlock := generatePowBlock(powChain[len(powChain)-1], m.PowID)

	if isPowBlockValid(newPowBlock, powChain[len(powChain)-1]) {
		powChain = append(powChain, newPowBlock)
		spew.Dump(powChain[len(powChain)-1])

		// After every committeeSize blocks, create a new committee and
		// add a block to comChain with new commmittee
		if len(powChain)%committeeSize == 1 {
			newComBlock := generateComBlock(len(powChain))
			comChain = append(comChain, newComBlock)
		}
	}
	mutex.Unlock()

	respondWithJSON(w, r, http.StatusCreated, newPowBlock)
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

// SHA256 hasing
func powBlockHash(block powBlock) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.PowID + block.PrevHash + block.Nonce
	return calculateHash(record)
}

func comBlockHash(block comBlock) string {
	committee := ""
	for i := 0; i < len(block.Committee); i++ {
		committee += block.Committee[i]
	}
	record := strconv.Itoa(block.Index) + block.Timestamp + committee + block.PrevHash
	return calculateHash(record)
}

func txBlockHash(block txBlock) string {
	txs := ""
	for i := 0; i < len(block.Transactions); i++ {
		txs += block.Transactions[i]
	}
	record := strconv.Itoa(block.Index) + block.Timestamp + txs + block.ComSign + block.PrevHash
	return calculateHash(record)
}

func calculateHash(record string) string {
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new powBlock using previous block's hash
func generatePowBlock(oldBlock powBlock, PowID string) powBlock {
	var newBlock powBlock

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.PowID = PowID
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = powDifficulty

	for i := 0; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(powBlockHash(newBlock), newBlock.Difficulty) {
			fmt.Println(powBlockHash(newBlock), " do more work!")
			// time.Sleep(100 * time.Millisecond)
			continue
		} else {
			fmt.Println(powBlockHash(newBlock), " work done!")
			newBlock.Hash = powBlockHash(newBlock)
			break
		}

	}
	return newBlock
}

func generateComBlock(height int) comBlock {
	var committee []string
	for i := height - committeeSize; i < height; i++ {
		committee = append(committee, powChain[i].PowID)
	}
	oldBlock := comChain[len(comChain)-1]
	block := comBlock{len(comChain), time.Now().String(), committee, "", oldBlock.Hash}
	block.Hash = comBlockHash(block)
	return block
}

func generateTxBlock() txBlock {
	oldBlock := txChain[len(txChain)-1]
	comSign := ""
	committee := comChain[len(comChain)-1].Committee
	for i := 0; i < len(committee); i++ {
		comSign += committee[i]
	}

	block := txBlock{len(txChain), time.Now().String(), pendingTxs, comSign, "", oldBlock.Hash}
	pendingTxs = []string{}
	block.Hash = txBlockHash(block)
	return block
}

func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}
