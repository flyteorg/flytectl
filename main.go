package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/flyteorg/flytectl/cmd"

	"github.com/flyteorg/flytectl/pkg/util"
	"github.com/flyteorg/flytestdlib/logger"
)

func main() {
	ctx := context.TODO()
	var wg sync.WaitGroup
	var err error
	wg.Add(2)
	var message string
	go func() {
		message, _ = util.DetectNewVersion(ctx)
		wg.Done()
	}()
	go func() {
		err = cmd.ExecuteCmd()
		if err != nil {
			logger.Error(ctx, err)
		}
		wg.Done()
	}()
	wg.Wait()
	if len(message) > 0 {
		fmt.Println(message)
	}
	if err != nil {
		os.Exit(1)
	}
}
