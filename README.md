# Simple Go Blockchain Implementation

## Introduction

This project implements a fully functional blockchain system in Go. It provides the structures and functions for creating, managing, and interacting with a blockchain, including a Proof of Work consensus mechanism, Merkle Tree for efficient transaction verification, and a Wallet system for managing user identities and transactions.

## Features

- Complete blockchain structure with blocks and transactions
- Proof-of-Work consensus mechanism
- Merkle Tree implementation for efficient transaction verification
- Wallet system for managing user identities and transactions
- Multi-wallet management system
- Transaction creation, signing, and verification
- UTXO (Unspent Transaction Output) model
- Persistent storage using BadgerDB
- Support for multiple nodes (with unique IDs)
- Genesis block creation
- Block mining
- Transaction and block serialization/deserialization
- Wallet persistence and recovery

## Installation

```bash
$ git clone https://github.com/mapfumo/golang-blockchain.git
$ cd golang-blockchain
$ go build -o blockchain-cli main.go
```

## Usage

```bash
$ ./blockchain-cli
Usage:
 getbalance -address ADDRESS - get the balance for an address
 createblockchain -address ADDRESS creates a blockchain and sends genesis reward to address
 printchain - Prints the blocks in the chain
 send -from FROM -to TO -amount AMOUNT -mine - Send amount of coins. Then -mine flag is set, mine off of this node
 createwallet - Creates a new Wallet
 listaddresses - Lists the addresses in our wallet file
 reindexutxo - Rebuilds the UTXO set
 startnode -miner ADDRESS - Start a node with ID specified in NODE_ID env. var. -miner enables mining
```

## BadgerDB

This blockchain implementation uses [BadgerDB](https://github.com/dgraph-io/badger), a fast, persistent key-value store written in Go. BadgerDB is chosen for its efficient performance and reliability in handling large volumes of data, making it suitable for blockchain storage needs.

### Key Features of BadgerDB

- **Fast Reads and Writes:** BadgerDB provides high-performance key-value access, crucial for the rapid querying and updating required by a blockchain.
- **Persistent Storage:** It ensures that blockchain data is stored reliably across restarts, using an append-only log and a memory table.
- **Low Overhead:** Designed with a focus on minimal memory overhead and efficient disk usage.
- **Transaction Support:** Supports ACID transactions to ensure data integrity during updates.

### How BadgerDB is Used in the Project

- **Blockchain Data Storage:** Blocks, transactions, and UTXO sets are stored in BadgerDB for efficient querying and persistence.
- **Reindexing:** The reindexutxo command rebuilds the UTXO set by scanning and reindexing data stored in BadgerDB.
- **State Management:** Maintains blockchain state and history through BadgerDB’s transaction and snapshot capabilities.

## Consensus

The consensus algorithm is a simple implementation of Proof of Work.

## Wallet System

The wallet system provides the following features:

- Key pair generation using Elliptic Curve Digital Signature Algorithm (ECDSA)
- Public address generation with version byte and checksum
- Address validation
- Wallet serialization and deserialization using gob encoding
- Multi-wallet management
- Wallet persistence and recovery from file storage

## Resources

- [Golang program to implement a Merkle tree](https://www.tutorialspoint.com/golang-program-to-implement-a-merkle-tree)
- [Base58 - An easy-to-share set of characters](https://learnmeabitcoin.com/technical/keys/base58/#:~:text=Base58%20is%20a%20user%2Dfriendly,private%20keys%2C%20and%20extended%20keys.)
- [Building a Blockchain in Golang](https://www.youtube.com/watch?v=yPS-hEyOfi4)
