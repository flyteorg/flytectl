package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/flyteorg/flytectl/pkg/docker"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/enescakir/emoji"
	sandboxConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/sandbox"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
)

const (
	startShort = "Start the flyte sandbox"
	startLong  = `
Start will run the flyte sandbox cluster inside a docker container and setup the config that is required 
::

 bin/flytectl sandbox start
	
Mount your flytesnacks repository code inside sandbox 
::

 bin/flytectl sandbox start --flytesnacks=$HOME/flyteorg/flytesnacks 
Usage
	`
)

var volumes = []mount.Mount{
	{
		Type:   mount.TypeBind,
		Source: f.FilePathJoin(f.UserHomeDir(), ".flyte"),
		Target: docker.K3sDir,
	},
}

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

func startSandboxCluster(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("%v Please Check your docker client %v \n", emoji.GrimacingFace, emoji.Whale)
		return err
	}
	reader, err := startSandbox(ctx, cli)
	if err != nil {
		return err
	}
	docker.WaitForSandbox(reader, docker.SuccessMessage)
	return nil
}

func startSandbox(ctx context.Context, cli docker.Docker) (*bufio.Scanner, error) {
	fmt.Printf("%v Bootstrapping a brand new flyte cluster... %v %v\n", emoji.FactoryWorker, emoji.Hammer, emoji.Wrench)
	if err := docker.SetupFlyteDir(); err != nil {
		return nil, err
	}

	if err := docker.GetFlyteSandboxConfig(); err != nil {
		return nil, err
	}

	if err := docker.RemoveSandbox(ctx, cli, os.Stdin); err != nil {
		return nil, err
	}

	if len(sandboxConfig.DefaultConfig.SnacksRepo) > 0 {
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: sandboxConfig.DefaultConfig.SnacksRepo,
			Target: docker.FlyteSnackDir,
		})
	}

	os.Setenv("KUBECONFIG", docker.Kubeconfig)
	os.Setenv("FLYTECTL_CONFIG", docker.FlytectlConfig)
	if err := docker.PullDockerImage(ctx, cli, docker.ImageName); err != nil {
		return nil, err
	}

	exposedPorts, portBindings, _ := docker.GetSandboxPorts()
	ID, err := docker.StartContainer(ctx, cli, volumes, exposedPorts, portBindings, docker.FlyteSandboxClusterName, docker.ImageName)
	if err != nil {
		fmt.Printf("%v Something went wrong: Failed to start Sandbox container %v, Please check your docker client and try again. \n", emoji.GrimacingFace, emoji.Whale)
		return nil, fmt.Errorf("error: %v", err)
	}

	_, errCh := docker.WatchError(ctx, cli, ID)
	reader, err := docker.ReadLogs(ctx, cli, ID)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}
	go func() {
		err := <-errCh
		if err != nil {
			fmt.Printf("err: %v", err)
			os.Exit(0)
		}
	}()

	return reader, nil
}
