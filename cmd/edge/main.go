package main

import (
	"os"

	"github.com/edgepayments/ept-cli/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
