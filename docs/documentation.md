# Documentation

Below is the UML diagram and documentation of this Proof-of-Concept implementation of ASHWAchain.

![UML](/UML.png)

# main.go

`main.go` contains the centralized server simulating the P2P network. It recieves requests from nodes to post blocks with the identity that the node owns and then creates a block in the Proof-of-Work blockchain. It also simulates the PBFT in Consensus Layer and creates committee block and transaction blocks and adds them to their respective blockchains.

Run the server by the following command
```bash
go run node.go
```

PoW, Committee, and Transaction blockchains are stored in lists namely `powChain`, `comChain`, and `txChain`, respectively.

## Configuration

- Set the `committeeSize` variable to determine the size of committee
    ```go
    const committeeSize = 5
    ```
- Set the `txBlockSize` to determine the size of transaction block
    ```go
    const txBlockSize = 10
    ```

## HTTP handlers
- Make a POST request to `localhost:8080/` to add a block to the PoW blockchain. Here the body must contain the ID as a string that needs to be added in the block as json data.
    ```json
    {
        "PowID": "<Block ID>"
    }
    ```
    - This request will create a new valid block with the given ID and append it to the PoW blockchain.
- Make a GET request to `localhost:8080/` to get the latest block in the PoW blockchain
- Make a GET request to `localhost:8080/committee` to view the Committee blockchain
- Make a GET request to `localhost:8080/tx/<transaction data>` to post a transaction. The url must contain the transaction data as a string.


## Functions

### handleGetPowChain
- This function is an HTTP request handler for the GET request to the url `localhost:8080/`.
- It displays the `powChain` on the webpage.


### handleWritePowBlock

- This function is an HTTP request handler for the POST request to the url `localhost:8080/`.
- It takes the JSON data in the HTTP request, generates a new block for the `powChain` with the `PowID` sent in the request data, and adds that block to the `powChain`.
- After every `committeeSize` number of blocks added to `powChain`, it takes the Identities stored in these blocks, creates a `comBlock`, and adds it to the `comChain`.

### handleTransactions

- This function is an HTTP request handler for the GET request to the url `localhost:8080/tx/{txData}`.
- It takes in the transaction data as a string parameter in the url and appends it to the `pendingTxs`. Once the `pendingTxs` list reaches the `txBlockSize`, it creates a `txBlock` with these transactions and adds it to the `txChain`.  
- It also displays the `txChain` on the webpage.

### handleGetComChain

- This function is an HTTP request handler for the GET request to the url `localhost:8080/committee`.
- It displays the `comChain` on the webpage.

### generatePowBlock

- This function is used by the `handleWritePowBlock` function to generate a new `powBlock`.
- It takes 2 arguments: the latest block in the `powChain` and the `PowID` for the new block that is being generated.
- It generates a valid `powBlock` with the `PowID` using Proof of work, i.e. calculating the `nonce` with a valid hash for the `powDifficulty`, and returns that block.


# node.py

`node.py` simulates a node running ASHWAchain for the GO server running from the `main.go` file.

Run the node executable with the only argument as the node's ID as follows
```bash
./node.py <Node ID>
```

- `ID` stores the node's identity, passed as the first system argument when starting the node. 

- Each node logs the added `PowID` in the file `<ID>.log` in the logs directory.

- Set the value of variables `miningTimeVariance` and `avgMiningTime` to simulate the proof-of-work mining. The values represent the time in seconds (s).
    ```py
    miningTimeVariance = 2
    avgMiningTime = 6
    ```

- The nodes store a list of IDs in `idList` to put in the PoW block. These IDs are generated for each node in the `IDs.csv` by running the `run.sh` file.



## BroadcastBlock

- This function runs till the node adds all the Identities stored in the `idList` to the `powChain`.
- To simulate PoW mining, `time.sleep` is used with an average mining time variable `avgMiningTime` and a mining time variance variable `miningTimeVariance`. 
- If the block is added, then node will try to add the next ID in the `idList`, else it will mine a new block again with the same ID.


# run.sh

- This is a bash script to run the program. It takes 2 arguments as follows.
    - Argument #1 is the number of nodes.
    - Argument #2 is the nunber of PowIDs that needs to be generated for each node.

- It will first run the `generateIDs.py` to create `IDs.csv`.
- Next, it will run a loop to start all the nodes with their IDs.


