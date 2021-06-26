package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/phayes/freeport"

	"github.com/docker/docker/api/types/mount"
	"github.com/enescakir/emoji"
	sandboxConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/sandbox"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/docker"
)

const (
	startShort = "Start the flyte sandbox"
	startLong  = `
Start will run the flyte sandbox cluster inside a docker container and setup the config that is required 
::

 bin/flytectl sandbox start
	
Mount your flytesnacks repository code inside sandbox 
::

 bin/flytectl sandbox start --sourcesPath=$HOME/flyteorg/flytesnacks 
Usage
	`
)

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

func startSandboxCluster(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	reader, err := startSandbox(ctx, cli, os.Stdin)
	if err != nil {
		return err
	}
	docker.WaitForSandbox(reader, docker.SuccessMessage, docker.Ports)
	fmt.Println()
	return nil
}

func startSandbox(ctx context.Context, cli docker.Docker, reader io.Reader) (*bufio.Scanner, error) {
	fmt.Printf("%v Bootstrapping a brand new flyte cluster... %v %v\n", emoji.FactoryWorker, emoji.Hammer, emoji.Wrench)

	ports, err := getPort(docker.Ports)
	if err != nil {
		return nil, err
	}
	docker.Ports = ports

	if err := docker.SetupFlyteDir(); err != nil {
		return nil, err
	}
	name := "flyte-sandbox"
	if len(sandboxConfig.DefaultConfig.Name) > 0 {
		name = sandboxConfig.DefaultConfig.Name
	}
	if err := docker.GetFlyteSandboxConfig(); err != nil {
		return nil, err
	}

	if err := docker.RemoveSandbox(ctx, cli, reader, name); err != nil {
		return nil, err
	}

	if len(sandboxConfig.DefaultConfig.SourcesPath) > 0 {
		docker.Volumes = append(docker.Volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: sandboxConfig.DefaultConfig.SourcesPath,
			Target: docker.FlyteSnackDir,
		})
	}

	os.Setenv("KUBECONFIG", docker.Kubeconfig)
	os.Setenv("FLYTECTL_CONFIG", docker.FlytectlConfig)
	if err := docker.PullDockerImage(ctx, cli, docker.ImageName); err != nil {
		return nil, err
	}

	exposedPorts, portBindings, _ := docker.GetSandboxPorts(docker.Ports)
	ID, err := docker.StartContainer(ctx, cli, docker.Volumes, exposedPorts, portBindings, name, docker.ImageName)
	if err != nil {
		fmt.Printf("%v Something went wrong: Failed to start Sandbox container %v, Please check your docker client and try again. \n", emoji.GrimacingFace, emoji.Whale)
		return nil, err
	}

	_, errCh := docker.WatchError(ctx, cli, ID)
	logReader, err := docker.ReadLogs(ctx, cli, ID)
	if err != nil {
		return nil, err
	}
	go func() {
		err := <-errCh
		if err != nil {
			fmt.Printf("err: %v", err)
			os.Exit(1)
		}
	}()

	return logReader, nil
}

func getPort(ports map[string]int) (map[string]int, error) {
	p, err := freeport.GetFreePorts(4)
	if err != nil {
		return map[string]int{}, err
	}
	i := 0
	for k := range ports {
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%v", ports[k]))
		if err != nil {
			ports[k] = p[i]
			i++
			continue
		}
		ln.Close()
	}
	return ports, nil
}
