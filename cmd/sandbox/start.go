package sandbox

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/enescakir/emoji"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
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

	if err := setupFlytectlConfig(); err != nil {
		return err
	}

	if container := getSandbox(cli); container != nil {
		if cmdUtil.AskForConfirmation("delete existing sandbox cluster", os.Stdin) {
			if err := teardownSandboxCluster(ctx, []string{}, cmdCtx); err != nil {
				return err
			}
		}
	}

	os.Setenv("KUBECONFIG", Kubeconfig)
	os.Setenv("FLYTECTL_CONFIG", FlytectlConfig)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Something goes wrong with container status", r)
		}
	}()

	ID, err := startContainer(cli)
	if err != nil {
		fmt.Println("Something goes wrong. We are not able to start sandbox container, Please check your docker client and try again ")
		return fmt.Errorf("error: %v", err)
	}

	go watchError(cli, ID)
	if err := readLogs(cli, ID); err != nil {
		return err
	}

	fmt.Printf("Add KUBECONFIG and FLYTECTL_CONFIG to your environment variable \n")
	fmt.Printf("export KUBECONFIG=%v \n", Kubeconfig)
	fmt.Printf("export FLYTECTL_CONFIG=%v \n", FlytectlConfig)
	return nil
}
