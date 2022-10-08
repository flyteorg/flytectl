package sandbox

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/enescakir/emoji"
	"github.com/flyteorg/flytectl/pkg/configutil"
	"github.com/flyteorg/flytectl/pkg/docker"
	"github.com/flyteorg/flytectl/pkg/k8s"
	"github.com/kataras/tablewriter"
)

func Teardown(ctx context.Context, cli docker.Docker, verbose bool) error {
	c, err := docker.GetSandbox(ctx, cli)
	if err != nil {
		return err
	}
	if c != nil {
		logCtx, cancel := context.WithCancel(context.Background())
		if verbose {
			go func(ctx context.Context) {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						printTeardownLogs(ctx, cli)
					}

				}
			}(logCtx)
		}
		if err := cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			cancel()
			return err
		}
		cancel()
	}
	if err := configutil.ConfigCleanup(); err != nil {
		fmt.Printf("Config cleanup failed. Which Failed due to %v \n ", err)
	}
	if err := removeSandboxKubeContext(); err != nil {
		fmt.Printf("Kubecontext cleanup failed. Which Failed due to %v \n ", err)
	}
	fmt.Printf("%v %v Sandbox cluster is removed successfully. \n", emoji.Broom, emoji.Broom)
	return nil
}

func printTeardownLogs(ctx context.Context, cli docker.Docker) {
	fmt.Print("\n---- Verbose Logs ----\n")
	var data = os.Stdout
	table := tablewriter.NewWriter(data)
	table.SetHeader([]string{"ID", "From", "Type", "Action"})
	table.SetRowLine(true)
	table.ClearRows()
	table.SetAutoWrapText(true)
	table.SetAutoFormatHeaders(true)
	msgs, _ := cli.Events(ctx, types.EventsOptions{
		Since: "1m",
	})
	table.Render()
	for m := range msgs {
		table.RenderRowOnce([]string{m.ID, m.From, m.Type, m.Action})
	}
}

func removeSandboxKubeContext() error {
	k8sCtxMgr := k8s.NewK8sContextManager()
	return k8sCtxMgr.RemoveContext(sandboxContextName)
}
