package sandbox

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/flyteorg/flytestdlib/logger"

	"github.com/docker/docker/api/types/mount"
	"github.com/flyteorg/flytectl/clierrors"
	"github.com/flyteorg/flytectl/pkg/k8s"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"

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
	progressBarMessage           = "Waiting for flyte deployment"
	k8sEndpoint                  = "https://127.0.0.1:30086"
	flyteMinimumVersionSupported = "v0.14.0"
	generatedManifest            = "/flyteorg/share/flyte_generated.yaml"
)

var (
	FlyteManifest = f.FilePathJoin(f.UserHomeDir(), ".flyte", "flyte_generated.yaml")
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

	time.Sleep(5 * time.Second)
	k8sClient, err := k8s.GetK8sClient(docker.Kubeconfig, k8sEndpoint)
	if err != nil {
		return err
	}

	total, countChan, errChK8s := monitorFlyteDeployment(ctx, k8sClient.AppsV1(), k8sClient.CoreV1().Nodes())

	go func() {
		err := <-errChK8s
		if err != nil {
			fmt.Printf("err: %v \n", err)
			os.Exit(1)
		}
	}()

	util.ProgressBarForFlyteStatus(total, countChan, progressBarMessage)

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
		} else {
			if err := downloadFlyteManifest(sandboxConfig.DefaultConfig.Version); err != nil {
				return nil, err
			}

			if vol, err := mountVolume(FlyteManifest, generatedManifest); err != nil {
				return nil, err
			} else if vol != nil {
				volumes = append(volumes, *vol)
			}
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

func watchFlyteDeploymentStatus(ctx context.Context, appsClient v1.AppsV1Interface, nodeClient corev1k8s.NodeInterface) (chan int64, chan error) {
	var count = make(chan int64)
	var lastCount int64 = 0
	var errorChan = make(chan error)

	go func() {
		for {
			ready, err := k8s.GetCountOfReadyDeployment(ctx, appsClient)
			if err != nil {
				errorChan <- err
			}
			count <- ready - lastCount
			lastCount = ready
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
	return count, errorChan
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

func monitorFlyteDeployment(ctx context.Context, appsClient v1.AppsV1Interface, nodeClient corev1k8s.NodeInterface) (int64, chan int64, chan error) {
	countChan := make(chan int64)
	chanErr := make(chan error)
	var err error
	defer func() {
		if err != nil {
			chanErr <- err
		}
	}()
	total, err := k8s.GetFlyteDeploymentCount(ctx, appsClient)
	if err != nil {
		return total, countChan, chanErr
	}

	countChan, chanErr = watchFlyteDeploymentStatus(ctx, appsClient, nodeClient)

	return total, countChan, chanErr
}

func downloadFlyteManifest(version string) error {
	release, err := util.CheckVersionExist(version, "flyte")
	if err != nil {
		return err
	}
	var manifestURL = ""
	for _, v := range release.Assets {
		if v.GetName() == "flyte_sandbox_manifest.yaml" {
			manifestURL = *v.BrowserDownloadURL
		}
	}
	if len(manifestURL) > 0 {
		response, err := http.Get(manifestURL)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		if response.StatusCode != 200 {
			return fmt.Errorf("someting goes wrong while downloading the flyte release")
		}
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}
		if err := util.WriteIntoFile(data, FlyteManifest); err != nil {
			return err
		}
	}
	return nil
}
