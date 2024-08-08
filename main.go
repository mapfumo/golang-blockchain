package main

import (
	"os"

	cli "github.com/mapfumo/golang-blockchain/cli"
)



func main() {
	// Ensure the program exits cleanly to avoid database corruption.
	defer os.Exit(0)


	cmd := cli.CommandLine{}
	cmd.Run()
}
