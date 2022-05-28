package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const N = 6

type Keypair struct {
	privkey *rsa.PrivateKey
	pubkey  *rsa.PublicKey
}

type KnownNode struct {
	nodeID int
	url    string
	pubkey *rsa.PublicKey
}

type MyNode struct {
	nodeID int
	url    string
	keys   Keypair
}

var KnownNodes []*KnownNode
var ClientNode *KnownNode
var Myself MyNode

func init() {
	generateKeyFiles()
	KnownNodes = make([]*KnownNode, N)
	setCommittee()
}

func setCommittee() {
	comm, _ := filepath.Abs("../committee_nodes.txt")
	comm_bytes, _ := ioutil.ReadFile(comm)
	committee := strings.Split(string(comm_bytes), "\n")

	nodeID := -1
	for i := 0; i <= N; i++ {
		if i == N {
			nodeID = 500
		} else {
			nodeID, _ = strconv.Atoi(committee[i])
		}

		pubFile, _ := filepath.Abs(fmt.Sprintf("../Keys/%d_pub", nodeID))
		pubfbytes, err := ioutil.ReadFile(pubFile)
		if err != nil {
			panic(err)
		}
		pubblock, _ := pem.Decode(pubfbytes)
		if pubblock == nil {
			panic(fmt.Errorf("parse block occured error"))
		}
		pubkey, err := x509.ParsePKIXPublicKey(pubblock.Bytes)

		ipaddrFile, _ := filepath.Abs(fmt.Sprintf("../Keys/%d_socket", nodeID))
		bytes, _ := ioutil.ReadFile(ipaddrFile)

		if i == N {
			ClientNode = &KnownNode{
				nodeID,
				string(bytes),
				pubkey.(*rsa.PublicKey),
			}
		} else {
			KnownNodes[i] = &KnownNode{
				nodeID,
				string(bytes),
				pubkey.(*rsa.PublicKey),
			}
		}
	}
}

func setMyNode(nodeID int) {
	privKey, pubKey, err := getKeyPairByFile(nodeID)
	if err != nil {
		panic(err)
	}
	ipaddrFile, _ := filepath.Abs(fmt.Sprintf("../Keys/%d_socket", nodeID))
	bytes, _ := ioutil.ReadFile(ipaddrFile)

	Myself = MyNode{
		nodeID,
		string(bytes),
		Keypair{
			privKey,
			pubKey,
		},
	}
}

func getKeyPairByFile(nodeID int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privFile, _ := filepath.Abs(fmt.Sprintf("../Keys/%d_priv", nodeID))
	pubFile, _ := filepath.Abs(fmt.Sprintf("../Keys/%d_pub", nodeID))
	fbytes, err := ioutil.ReadFile(privFile)
	if err != nil {
		return nil, nil, err
	}
	block, _ := pem.Decode(fbytes)
	if block == nil {
		return nil, nil, fmt.Errorf("parse block occured error")
	}
	privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}
	pubfbytes, err := ioutil.ReadFile(pubFile)
	if err != nil {
		return nil, nil, err
	}
	pubblock, _ := pem.Decode(pubfbytes)
	if pubblock == nil {
		return nil, nil, fmt.Errorf("parse block occured error")
	}
	pubkey, err := x509.ParsePKIXPublicKey(pubblock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	return privkey, pubkey.(*rsa.PublicKey), nil
}

func generateKeyFiles() {
	if !FileExists("../Keys") {
		err := os.Mkdir("Keys", 0700)
		if err != nil {
			panic(err)
		}
		for i := 0; i <= N; i++ {
			filename, _ := filepath.Abs(fmt.Sprintf("../Keys/%d", i))
			if !FileExists(filename+"_priv") && !FileExists(filename+"_pub") {
				priv, pub := generateKeyPair()
				err := ioutil.WriteFile(filename+"_priv", priv, 0644)
				if err != nil {
					panic(err)
				}
				ioutil.WriteFile(filename+"_pub", pub, 0644)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func generateKeyPair() ([]byte, []byte) {
	privkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	mprivkey := x509.MarshalPKCS1PrivateKey(privkey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: mprivkey,
	}
	bprivkey := pem.EncodeToMemory(block)
	pubkey := &privkey.PublicKey
	mpubkey, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: mpubkey,
	}
	bpubkey := pem.EncodeToMemory(block)
	return bprivkey, bpubkey
}
