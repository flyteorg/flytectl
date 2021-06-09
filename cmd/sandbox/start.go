package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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
	prompt      = promptui.Prompt{
		Label:     "Delete Existing Sandbox Cluster",
		IsConfirm: true,
	}
)

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

func startSandboxCluster(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {

	fmt.Printf("%v It will take some time, We will start a fresh flyte cluster for you %v %v\n", emoji.ManTechnologist, emoji.Rocket, emoji.Rocket)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("Please Check your docker client %v \n", emoji.ManTechnologist)
		return err
	}

	if container := getSandbox(cli); container != nil {
		answer, err := prompt.Run()
		if err != nil {
			return err
		}
		if strings.ToLower(answer) == "y" {
			if err := teardownSandboxCluster(ctx, []string{}, cmdCtx); err != nil {
				return err
			}
		}
	}
	if err := setupFlytectlConfig(); err != nil {
		return err
	}

	ID, err := startContainer(cli)
	if err == nil {
		if err := cli.ContainerStart(ctx, ID, types.ContainerStartOptions{}); err != nil {
			return err
		}

		os.Setenv("KUBECONFIG", KUBECONFIG)

		go func() {
			statusCh, errCh := cli.ContainerWait(ctx, ID, container.WaitConditionNotRunning)
			select {
			case err := <-errCh:
				if err != nil {
					panic(err)
				}
			case <-statusCh:
			}
		}()

		reader, _ := cli.ContainerLogs(context.Background(), ID, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		})

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
		fmt.Printf("Add (KUBECONFIG) to your ENV variabl \n")
		fmt.Printf("export KUBECONFIG=%v \n", KUBECONFIG)
		return nil
	}
	fmt.Println("Something goes wrong. We are not able to start sandbox container, Please check your docker client and try again \n", emoji.Rocket)
	fmt.Printf("error: %v", err)
	return nil
}
