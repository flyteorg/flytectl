package demo

import (
	"context"
	"fmt"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/docker"
)

const (
	reloadShort = "Power cycle the Flyte executable pod, effectively picking up the "
	reloadLong  = `
If you've changed the ~/.flyte/state/flyte.yaml file, run this command to restart the Flyte binary pod, effectively
picking up the new settings:

Usage
::

 flytectl demo reload

`
)

func reloadDemoCluster(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}
	fmt.Println(cli)
	// Need to access the k3s cluster running inside the flyte-sandbox docker container.
	// Look for a pod by name,
	// Power
	return nil
}
