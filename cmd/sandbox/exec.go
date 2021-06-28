package sandbox

import (
	"context"
	"fmt"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/docker"
)

const (
	execShort = "Get the status of the sandbox environment."
	execLong  = `
Status will retrieve the status of the Sandbox environment. Currently FlyteSandbox runs as a local docker container.
This will return the docker status for this container
Usage
::
 bin/flytectl sandbox status 
`
)

func sandboxClusterExec(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}
	if len(args) > 0 {
		return Execute(ctx, cli, args)
	}
	return fmt.Errorf("Please use the right syntax. flytectl sandbox exec -- kubectl get pods")
}

func Execute(ctx context.Context, cli docker.Docker, args []string) error {
	c := docker.GetSandbox(ctx, cli)
	if c != nil {
		exec, err := docker.ExecCommend(ctx, cli, c.ID, args)
		if err != nil {
			return err
		}
		if err := docker.InspectExecResp(ctx, cli, exec.ID); err != nil {
			return err
		}
	}
	return nil
}
