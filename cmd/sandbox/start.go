package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/flyteorg/flytectl/clierrors"
	"github.com/flyteorg/flytectl/pkg/util/githubutil"

	"github.com/avast/retry-go"
	"github.com/olekukonko/tablewriter"
	corev1api "k8s.io/api/core/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/docker/docker/api/types/mount"
	"github.com/flyteorg/flytectl/pkg/configutil"
	"github.com/flyteorg/flytectl/pkg/k8s"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/enescakir/emoji"
	sandboxConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/sandbox"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/docker"
	"github.com/flyteorg/flytectl/pkg/util"
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
	
Run specific version of flyte. flytectl sandbox only support flyte version available in Github release https://github.com/flyteorg/flyte/tags
::

 bin/flytectl sandbox start  --version=v0.14.0

Note: Flytectl sandbox will only work for v0.10.0+
	
Usage
	`
	k8sEndpoint             = "https://127.0.0.1:30086"
	flyteNamespace          = "flyte"
	flyteRepository         = "flyte"
	dind                    = "dind"
	sandboxSupportedVersion = "v0.10.0"
	diskPressureTaint       = "node.kubernetes.io/disk-pressure"
	taintEffect             = "NoSchedule"
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

	if reader != nil {
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
		if err := watchFlyteDeployment(ctx, k8sClient.CoreV1()); err != nil {
			return err
		}
		util.PrintSandboxMessage()
	}
	return nil
}

func startSandbox(ctx context.Context, cli docker.Docker, reader io.Reader) (*bufio.Scanner, error) {
	fmt.Printf("%v Bootstrapping a brand new flyte cluster... %v %v\n", emoji.FactoryWorker, emoji.Hammer, emoji.Wrench)

	if err := docker.RemoveSandbox(ctx, cli, reader); err != nil {
		if err.Error() != clierrors.ErrSandboxExists {
			return nil, err
		}
		fmt.Printf("Existing details of your sandbox")
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

	var tag = "dind"
	if len(sandboxConfig.DefaultConfig.Version) > 0 {
		isGreater, err := util.IsVersionGreaterThan(sandboxConfig.DefaultConfig.Version, sandboxSupportedVersion)
		if err != nil {
			return nil, err
		}
		if !isGreater {
			return nil, fmt.Errorf("version flag only supported with flyte %s+ release", sandboxSupportedVersion)
		}
		sha, err := githubutil.GetSHAFromVersion(sandboxConfig.DefaultConfig.Version, flyteRepository)
		if err != nil {
			return nil, err
		}
		tag = fmt.Sprintf("%s-%s", dind, sha)
	}

	// Latest release will use image cr.flyte.org/flyteorg/flyte-sandbox:dind
	// In case of version flytectl will use cr.flyte.org/flyteorg/flyte-sandbox:dind-{SHA}
	var sandboxImageName = fmt.Sprintf("%s:%s", docker.ImageName, tag)

	fmt.Printf("%v pulling docker image for release %s\n", emoji.Whale, sandboxImageName)
	if err := docker.PullDockerImage(ctx, cli, sandboxImageName); err != nil {
		return nil, err
	}

	fmt.Printf("%v booting Flyte-sandbox container\n", emoji.FactoryWorker)
	exposedPorts, portBindings, _ := docker.GetSandboxPorts()
	ID, err := docker.StartContainer(ctx, cli, volumes, exposedPorts, portBindings, docker.FlyteSandboxClusterName, sandboxImageName)
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

func watchFlyteDeployment(ctx context.Context, appsClient corev1.CoreV1Interface) error {
	var data = os.Stdout
	table := tablewriter.NewWriter(data)
	table.SetHeader([]string{"Service", "Status", "Namespace"})
	table.SetRowLine(true)

	for {
		isTaint, err := isNodeTainted(ctx, appsClient)
		if err != nil {
			return err
		}
		if isTaint {
			return fmt.Errorf("docker sandbox doesn't have sufficient memory available. Please run docker system prune -a --volumes")
		}

		pods, err := getFlyteDeployment(ctx, appsClient)
		if err != nil {
			return err
		}
		table.ClearRows()
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)

		// Clear os.Stdout
		_, _ = data.WriteString("\x1b[3;J\x1b[H\x1b[2J")

		var total, ready int
		total = len(pods.Items)
		ready = 0
		if total != 0 {
			for _, v := range pods.Items {
				if isPodReady(v) {
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

		time.Sleep(40 * time.Second)
	}

	return nil
}

func isPodReady(v corev1api.Pod) bool {
	if (v.Status.Phase == corev1api.PodRunning) || (v.Status.Phase == corev1api.PodSucceeded) {
		return true
	}
	return false
}

func getFlyteDeployment(ctx context.Context, client corev1.CoreV1Interface) (*corev1api.PodList, error) {
	pods, err := client.Pods(flyteNamespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func isNodeTainted(ctx context.Context, client corev1.CoreV1Interface) (bool, error) {
	nodes, err := client.Nodes().List(ctx, v1.ListOptions{})
	if err != nil {
		return false, err
	}
	match := 0
	for _, node := range nodes.Items {
		for _, c := range node.Spec.Taints {
			if c.Key == diskPressureTaint && c.Effect == taintEffect {
				match++
			}
		}
	}
	if match > 0 {
		return true, nil
	}
	return false, nil
}
