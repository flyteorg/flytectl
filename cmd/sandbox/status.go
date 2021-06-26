package sandbox

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
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
	messageTemplate = `Flyte UI is available at http://localhost:%v/console
Kubernetes Dashboard is available at http://localhost:%v
Minio Dashboard is available at http://localhost:%v
Flyteadmin endpoint is available at http://localhost:%v`
)

func sandboxClusterStatus(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	return printStatus(ctx, cli, args)
}

func printStatus(ctx context.Context, cli docker.Docker, args []string) error {
	if len(args) > 0 {
		c := docker.GetSandbox(ctx, cli, args[0])
		if c != nil {
			fmt.Printf("Flyte sadasd sandbox cluster[%v] container image [%s] with status [%s] is in state [%s] \n\n", strings.TrimPrefix(c.Names[0], "/"), c.Image, c.Status, c.State)
			for _, v := range c.Ports {
				k := printPortInformation(v)
				docker.Ports[k] = int(v.PublicPort)
			}
			fmt.Printf(messageTemplate, docker.Ports["console"], docker.Ports["k8s"], docker.Ports["minio"], docker.Ports["admin"])
			return nil
		}
		fmt.Printf("No sandbox cluster found with name [%v] \n", args[0])
		return nil
	}
	containers := docker.GetAllSandbox(ctx, cli)
	for _, c := range containers {
		fmt.Printf("Flyte local sandbox cluster[%v] container image [%s] with status [%s] is in state [%s] \n", strings.TrimPrefix(c.Names[0], "/"), c.Image, c.Status, c.State)
	}
	return nil
}

func printPortInformation(v types.Port) string {
	switch v.PrivatePort {
	case 30081:
		return "k8s"
	case 30084:
		return "k8s"
	case 30086:
		return "minio"
	case 30082:
		return "admin"
	}
	return ""
}
