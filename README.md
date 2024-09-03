# Simple Blockchain in Golang

## Introduction

This is a simple blockchain implementation in Golang.

## Consensus

The consensus algorithm is a simple implementation of Proof of Work.

## Notes

[Big Int in Go](https://pkg.go.dev/math/big) package used this app consensus algorithm provides support for large numbers (with arbitrary precision) and includes types like Int for big integers, Rat for rational numbers, and Float for floating-point numbers.

[_badger_](https://pkg.go.dev/github.com/dgraph-io/badger/v4) is a fast, open source, key-value database written in Go. This is where the blockchain data is stored. Badger db only accepts arrays of bytes or slices of bytes we need to seriialise our blockchain data structures into bytes

reindexutxo - Rebuilds the UTXO set

# go run man.go reindexutxo
