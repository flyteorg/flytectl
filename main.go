package main

import (
	"os"

	"github.com/flyteorg/flytectl/cmd"
)

func main() {
	if err := cmd.ExecuteCmd(); err != nil {
		os.Exit(1)
	}
}
