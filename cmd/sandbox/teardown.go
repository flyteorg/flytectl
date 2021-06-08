package sandbox

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	"os"
)


const (
	teardownShort = "Teardown will cleanup the sandbox environment"
	teardownLong  = `
Teardown will remove docker container and all the flyte config 
::

 bin/flytectl teardown 


Usage
`
)

func teardownSandboxCluster(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	for _, container := range containers {
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			return err
		}
		if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			return err
		}
	}

	err = os.Remove(f.FilePathJoin(f.UserHomeDir(), ".flyte","config.yaml"))
	if err != nil {
		return err
	}
	err = os.Remove(f.FilePathJoin(f.UserHomeDir(), ".flyte","kube.yaml"))
	if err != nil {
		return err
	}
	return  nil
}
