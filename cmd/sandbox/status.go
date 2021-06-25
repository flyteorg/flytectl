package sandbox

import (
	"context"
	"fmt"

	sandboxConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/sandbox"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/docker"
)

const (
	statusShort = "Get the status of the sandbox environment."
	statusLong  = `
Status will retrieve the status of the Sandbox environment. Currently FlyteSandbox runs as a local docker container.
This will return the docker status for this container

Usage
::

 bin/flytectl sandbox status 

`
)

func sandboxClusterStatus(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	return printStatus(ctx, cli)
}

func printStatus(ctx context.Context, cli docker.Docker) error {
	containers := docker.GetAllSandbox(ctx, cli, sandboxConfig.DefaultConfig.Name)
	for _, c := range containers {
		fmt.Printf("Flyte local sandbox cluster container image [%s] with status [%s] is in state [%s]", c.Image, c.Status, c.State)
		return nil
	}
	return nil
}
