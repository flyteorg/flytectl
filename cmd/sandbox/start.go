package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/avast/retry-go"
	"github.com/olekukonko/tablewriter"
	corev1api "k8s.io/api/core/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/flyteorg/flytectl/pkg/util/githubutil"

	"github.com/flyteorg/flytestdlib/logger"

	"github.com/docker/docker/api/types/mount"
	"github.com/flyteorg/flytectl/clierrors"
	"github.com/flyteorg/flytectl/pkg/k8s"

	"github.com/flyteorg/flytectl/pkg/configutil"

	"github.com/flyteorg/flytectl/pkg/docker"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	"github.com/flyteorg/flytectl/pkg/util"
	corev1k8s "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/enescakir/emoji"
	sandboxConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/sandbox"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
)

const (
	startShort = "Start the flyte sandbox cluster"
	startLong  = `
The Flyte Sandbox is a fully standalone minimal environment for running Flyte. provides a simplified way of running flyte-sandbox as a single Docker container running locally.  

Start sandbox cluster without any source code
::

 bin/flytectl sandbox start
	
Mount your source code repository inside sandbox 
::

 bin/flytectl sandbox start --source=$HOME/flyteorg/flytesnacks 
	
Run specific version of flyte, Only available after v0.14.0+
::

 bin/flytectl sandbox start  --version=v0.14.0

Usage
	`
	k8sEndpoint                  = "https://127.0.0.1:30086"
	flyteMinimumVersionSupported = "v0.14.0"
	generatedManifest            = "/flyteorg/share/flyte_generated.yaml"
)

var (
	flyteManifest = f.FilePathJoin(f.UserHomeDir(), ".flyte", "flyte_generated.yaml")
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
	if reader != nil {
		docker.WaitForSandbox(reader, docker.SuccessMessage)
	}

	var k8sClient k8s.K8s
	err = retry.Do(
		func() error {
			k8sClient, err = k8s.GetK8sClient(docker.Kubeconfig, k8sEndpoint)
			return err
		},
		retry.Attempts(10),
	)
	if err != nil {
		return err
	}

	podsChan, errChK8s := watchFlyteDeployment(ctx, k8sClient.CoreV1())
	go func() {
		err := <-watchDiskPressure(ctx, k8sClient.CoreV1().Nodes(), errChK8s)
		if err != nil {
			fmt.Printf("err: %v \n", err)
			os.Exit(1)
		}
	}()

	var data = os.Stdout
	table := tablewriter.NewWriter(data)
	table.SetHeader([]string{"Service", "Status", "Namespace"})
	table.SetRowLine(true)

	var total, ready int
	for pods := range podsChan {
		table.ClearRows()
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		_, _ = data.WriteString("\x1b[3;J\x1b[H\x1b[2J")
		total = len(pods.Items)
		ready = 0
		if total != 0 {
			for _, v := range pods.Items {
				if (v.Status.Phase == corev1api.PodRunning) || (v.Status.Phase == corev1api.PodSucceeded) {
					ready++
				}
				if len(v.Status.Conditions) > 0 {
					table.Append([]string{v.GetName(), string(v.Status.Phase), v.GetNamespace()})
				}
			}
			table.Render()
			if total == ready {
				break
			}
		}

	}
	util.PrintSandboxMessage()
	return nil
}

func startSandbox(ctx context.Context, cli docker.Docker, reader io.Reader) (*bufio.Scanner, error) {
	fmt.Printf("%v Bootstrapping a brand new flyte cluster... %v %v\n", emoji.FactoryWorker, emoji.Hammer, emoji.Wrench)

	if err := docker.RemoveSandbox(ctx, cli, reader); err != nil {
		if err.Error() != clierrors.ErrSandboxExists {
			return nil, err
		}
		fmt.Printf("Existing details of your sandbox:")
		util.PrintSandboxMessage()
		return nil, nil
	}

	if err := util.SetupFlyteDir(); err != nil {
		return nil, err
	}

	templateValues := configutil.ConfigTemplateSpec{
		Host:     "localhost:30081",
		Insecure: true,
	}
	if err := configutil.SetupConfig(configutil.FlytectlConfig, configutil.GetSandboxTemplate(), templateValues); err != nil {
		return nil, err
	}

	volumes := docker.Volumes
	if vol, err := mountVolume(sandboxConfig.DefaultConfig.Source, docker.Source); err != nil {
		return nil, err
	} else if vol != nil {
		volumes = append(volumes, *vol)
	}

	if len(sandboxConfig.DefaultConfig.Version) > 0 {
		isGreater, err := util.IsVersionGreaterThan(sandboxConfig.DefaultConfig.Version, flyteMinimumVersionSupported)
		if err != nil {
			return nil, err
		}
		if !isGreater {
			logger.Infof(ctx, "version flag only supported after with flyte %s+ release", flyteMinimumVersionSupported)
			return nil, fmt.Errorf("version flag only supported after with flyte %s+ release", flyteMinimumVersionSupported)
		}
		if err := githubutil.GetFlyteManifest(sandboxConfig.DefaultConfig.Version, flyteManifest); err != nil {
			return nil, err
		}

		if vol, err := mountVolume(flyteManifest, generatedManifest); err != nil {
			return nil, err
		} else if vol != nil {
			volumes = append(volumes, *vol)
		}

	}

	fmt.Printf("%v pulling docker image %s\n", emoji.Whale, docker.ImageName)
	if err := docker.PullDockerImage(ctx, cli, docker.ImageName); err != nil {
		return nil, err
	}

	fmt.Printf("%v booting Flyte-sandbox container\n", emoji.FactoryWorker)
	exposedPorts, portBindings, _ := docker.GetSandboxPorts()
	ID, err := docker.StartContainer(ctx, cli, volumes, exposedPorts, portBindings, docker.FlyteSandboxClusterName, docker.ImageName)
	if err != nil {
		fmt.Printf("%v Something went wrong: Failed to start Sandbox container %v, Please check your docker client and try again. \n", emoji.GrimacingFace, emoji.Whale)
		return nil, err
	}

	logReader, err := docker.ReadLogs(ctx, cli, ID)
	if err != nil {
		return nil, err
	}

	return logReader, nil
}

func watchDiskPressure(ctx context.Context, nodeClient corev1k8s.NodeInterface, errorChan chan error) chan error {
	go func() {
		for {
			isTaint, err := k8s.GetNodeTaintStatus(ctx, nodeClient)
			if err != nil {
				errorChan <- err
			}
			if isTaint {
				errorChan <- fmt.Errorf("docker sandbox doesn't have sufficient memory available. Please run docker system prune -a --volumes")
			}
			time.Sleep(5 * time.Second)
		}
	}()
	return errorChan
}

func mountVolume(file, destination string) (*mount.Mount, error) {
	if len(file) > 0 {
		source, err := filepath.Abs(file)
		if err != nil {
			return nil, err
		}
		return &mount.Mount{
			Type:   mount.TypeBind,
			Source: source,
			Target: destination,
		}, nil
	}
	return nil, nil
}

func watchFlyteDeployment(ctx context.Context, appsClient corev1.CoreV1Interface) (chan corev1api.PodList, chan error) {
	watchPods := make(chan corev1api.PodList)
	chanErr := make(chan error)

	go func() {
		for {
			pods, err := k8s.GetFlyteDeployment(ctx, appsClient)
			if err != nil {
				chanErr <- err
			}
			watchPods <- *pods
			time.Sleep(30 * time.Second)
		}
	}()

	return watchPods, chanErr
}
