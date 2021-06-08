package sandbox

import (
	"context"
	"fmt"
	"strings"

	"github.com/enescakir/emoji"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
)

const (
	teardownShort = "Teardown will cleanup the sandbox environment"
	teardownLong  = `
Teardown will remove docker container and all the flyte config 
::

 bin/flytectl sandbox teardown 

Stop will remove docker container and all the flyte config 
::

 bin/flytectl sandbox stop 


Usage
`
)

func teardownSandboxCluster(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return err
	}
	var containerID string
	var isExist = false
	for _, container := range containers {
		if strings.Contains(container.Names[0], SandboxClusterName) {
			containerID = container.ID
			isExist = true
		}
	}
	if isExist {
		if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			return err
		}
	}

	if err := ConfigCleanup(); err != nil {
		fmt.Printf("Config cleanup failed. You can manually remove them from ~/.flyte directory. %v \n ", err)
	}
	fmt.Printf("Sandbox cluster is removed successfully %v \n", emoji.Rocket)
	return nil
}
