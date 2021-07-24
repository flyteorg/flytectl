package main

import (
	"context"
	"os"

	"github.com/flyteorg/flytectl/pkg/util"

	"github.com/flyteorg/flytectl/cmd"
	"github.com/flyteorg/flytestdlib/logger"
)

func main() {
	ctx := context.TODO()
	if err := cmd.ExecuteCmd(); err != nil {
		logger.Error(ctx, err)
		os.Exit(1)
	}
	util.DetectNewVersion(ctx)
}
