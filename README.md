# ASHWAchain

A Proof-of-Concept implementation of ASHWAchain protocol.

# Prerequisites

- go version go1.10.4 linux/amd64
- Python 3.6.9

# How to run
- Open a terminal, navigate to the root directory and run the go server
    ```bash 
    go run main.go
    ```

- Provide executable persmissions to all `.py` and `.sh` files
    ```
    chmod +x *.py *.sh
    ```

- Open another terminal and run the executable with 2 arguments
    - Argument 1: Total nodes
    - Arguemnt 2: Number of identities for each node to use
    ```
    ./run.sh <# NODES> <# IDENTITIES PER NODE>
    ```

- Go to 'localhost:8080/' to view the Proof-of-Work blockchain

## Create Transactions
- Make a GET request with transaction data as a string in the URL as follows
    - `localhost:8080/tx/<Transaction Data>`

## Stop the Nodes and clear Logs
- Install `killall` using `sudo apt-get install psmisc`
- Run the `./clear.sh`


# Future work
- Migrate the generation of PoW block in the nodes instead of the GO server.
- Create a real P2P network.
- Add operation log, and other states like account balances, etc.

