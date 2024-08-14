package main

import (
	"os"

	"github.com/mapfumo/golang-blockchain/wallet"
)

func main() {
	// Ensure the program exits cleanly to avoid database corruption.
	defer os.Exit(0)

	// cmd := cli.CommandLine{}
	// cmd.Run()

	w := wallet.MakeWallet()
	w.Address()
}
