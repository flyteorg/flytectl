package main

import "github.com/flyteorg/flytectl/cmd"

func main() {
	if err := cmd.ExecuteCmd(); err != nil {
		panic(err)
	}
}
