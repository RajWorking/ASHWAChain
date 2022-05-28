package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Client struct {
	nodeId     int
	url        string
	keypair    Keypair
	knownNodes []*KnownNode
	request    *RequestMsg
	replyLog   map[int]*ReplyMsg
	mutex      sync.Mutex
}

func NewClient() *Client {
	setMyNode(ClientNode.nodeID)
	rand.Seed(time.Now().UnixNano())
	client := &Client{
		ClientNode.nodeID,
		ClientNode.url,
		Myself.keys,
		KnownNodes,
		nil,
		make(map[int]*ReplyMsg),
		sync.Mutex{},
	}
	return client
}

func (c *Client) Start() {
	c.sendRequest()
	fmt.Println(c.url)
	ln, err := net.Listen("tcp", c.url)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for i := 0; i < N; i++ {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go c.handleConnection(conn)

		if c.checkSuccess() {
			return
		}
	}

	for {
		if c.checkSuccess() {
			break
		}
	}
}

func (c *Client) checkSuccess() bool {
	c.mutex.Lock()
	rlen := len(c.replyLog)
	c.mutex.Unlock()

	fmt.Println("rlen: ", rlen, "of", c.countNeedReceiveMsgAmount())

	if rlen >= c.countNeedReceiveMsgAmount() {
		fmt.Println("request success!!")
		return true
	}

	return false
}

func (c *Client) handleConnection(conn net.Conn) {
	req, err := ioutil.ReadAll(conn)
	header, payload, _ := SplitMsg(req)
	if err != nil {
		panic(err)
	}
	switch header {
	case hReply:
		c.handleReply(payload)
	}
}

func (c *Client) sendRequest() {
	bytes, _ := ioutil.ReadFile("../transaction.txt")
	msg := string(bytes)

	req := Request{
		msg,
		hex.EncodeToString(generateDigest(msg)),
	}
	reqmsg := &RequestMsg{
		"validate",
		int(time.Now().Unix()),
		c.nodeId,
		req,
	}
	sig, err := c.signMessage(reqmsg)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	logBroadcastMsg(hRequest, reqmsg)
	send(ComposeMsg(hRequest, reqmsg, sig), c.findPrimaryNode().url)
	c.request = reqmsg
}

func (c *Client) handleReply(payload []byte) {
	var replyMsg ReplyMsg
	err := json.Unmarshal(payload, &replyMsg)
	if err != nil {
		fmt.Printf("error happened:%v", err)
		return
	}
	logHandleMsg(hReply, replyMsg, replyMsg.NodeID)
	c.mutex.Lock()
	c.replyLog[replyMsg.NodeID] = &replyMsg
	c.mutex.Unlock()
}

func (c *Client) signMessage(msg interface{}) ([]byte, error) {
	sig, err := signMessage(msg, c.keypair.privkey)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (c *Client) findPrimaryNode() *KnownNode {
	nodeId := c.knownNodes[ViewID%len(c.knownNodes)].nodeID
	for _, knownNode := range c.knownNodes {
		if knownNode.nodeID == nodeId {
			return knownNode
		}
	}
	return nil
}

func (c *Client) countNeedReceiveMsgAmount() int {
	f := (len(c.knownNodes) - 1) / 3
	return f + 1
}
