package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
)

var (
	Environment = []string{"SANDBOX=1", "KUBERNETES_API_PORT=30086","FLYTE_HOST=localhost:30081","FLYTE_AWS_ENDPOINT=http://localhost:30084"}
)


func startSandboxCluster(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {

	response, err := http.Get("https://raw.githubusercontent.com/flyteorg/flytectl/master/config.yaml")
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.FilePathJoin(f.UserHomeDir(), ".flyte","config.yaml"), data , 0644)
	if err != nil {
		return err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	r, err := cli.ImagePull(ctx, ImageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, r)

	exposedPorts, portBindings, _ := nat.ParsePortSpecs([]string{
		"127.0.0.1:30086:30086",
		"127.0.0.1:30081:30081",
		"127.0.0.1:30082:30082",
		"127.0.0.1:30084:30084",
	})

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Env: Environment,
		Image: ImageName,
		Tty:   false,
		ExposedPorts: exposedPorts,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: f.FilePathJoin(f.UserHomeDir(), "kubeconfig"),
				Target: "/etc/rancher/",
			},
		},
		PortBindings: portBindings,
		Privileged: true,
	}, nil,
		nil,"flyte-sandbox")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	kubeconfig, err := ioutil.ReadFile(f.FilePathJoin(f.UserHomeDir(), "kubeconfig","k3s","k3s.yaml"))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(f.FilePathJoin(f.UserHomeDir(), ".flyte","kube.yaml"), kubeconfig, 0644)
		if err != nil {
			return err
		}

	os.Setenv("KUBECONFIG",f.FilePathJoin(f.UserHomeDir(), ".flyte","kube.yaml"))

	reader, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
		Follow:     true,
	})
	defer reader.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return  nil
}
