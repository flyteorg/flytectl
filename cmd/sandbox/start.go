package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/go-connections/nat"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/enescakir/emoji"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
)

const (
	startShort = "Start the flyte sandbox"
	startLong  = `
Start will run the flyte sandbox cluster inside a docker container and setup the config that is required 
::

 bin/flytectl start

Usage
	`
	ImageName = "ghcr.io/flyteorg/flyte-sandbox:dind"

	FlytectlConfig     = "https://raw.githubusercontent.com/flyteorg/flytectl/master/config.yaml"
	SandboxClusterName = "flyte-sandbox"
)

var (
	Environment = []string{"SANDBOX=1", "KUBERNETES_API_PORT=30086", "FLYTE_HOST=localhost:30081", "FLYTE_AWS_ENDPOINT=http://localhost:30084"}
)

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

func startSandboxCluster(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {

	fmt.Printf("%v It will take some time, We will start a fresh flyte cluster for you %v %v\n", emoji.ManTechnologist, emoji.Rocket, emoji.Rocket)

	if err := SetupFlytectlConfig(); err != nil {
		return err
	}
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("Please Check your docker client %v \n", emoji.ManTechnologist)
		return err
	}
	r, err := cli.ImagePull(ctx, ImageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	if _, err := io.Copy(os.Stdout, r); err != nil {
		return err
	}

	ExposedPorts, PortBindings, _ := nat.ParsePortSpecs([]string{
		"127.0.0.1:30086:30086",
		"127.0.0.1:30081:30081",
		"127.0.0.1:30082:30082",
		"127.0.0.1:30084:30084",
	})

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Env:          Environment,
		Image:        ImageName,
		Tty:          false,
		ExposedPorts: ExposedPorts,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: f.FilePathJoin(f.UserHomeDir(), ".flyte"),
				Target: "/etc/rancher/",
			},
			// TODO (Yuvraj) Add flytectl config in sandbox and mount with host file system
			//{
			//	Type:   mount.TypeBind,
			//	Source: f.FilePathJoin(f.UserHomeDir(), ".flyte", "config.yaml"),
			//	Target: "/.flyte/",
			//},
		},
		PortBindings: PortBindings,
		Privileged:   true,
	}, nil,
		nil, SandboxClusterName)
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	os.Setenv("KUBECONFIG", KUBECONFIG)

	reader, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true,
	})
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "Flyte is ready! Flyte UI is available at http://localhost:30081/console") {
			fmt.Printf("%v Flyte is ready! Flyte UI is available at http://localhost:30081/console. %v %v %v \n", emoji.ManTechnologist, emoji.Rocket, emoji.Rocket, emoji.PartyPopper)
			fmt.Printf("Please visit https://github.com/flyteorg/flytesnacks for more example %v \n", emoji.Rocket)
			fmt.Printf("Register all flytesnacks example by running 'flytectl register examples  -d development  -p flytesnacks' \n")
			break
		}
		fmt.Println(scanner.Text())
	}
	return nil
}
