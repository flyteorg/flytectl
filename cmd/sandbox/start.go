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
	"strings"

	"github.com/flyteorg/flytestdlib/logger"

	"github.com/flyteorg/flytectl/clierrors"

	"github.com/docker/docker/api/types/mount"

	"github.com/flyteorg/flytectl/pkg/configutil"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	"github.com/flyteorg/flytectl/pkg/util"

	"github.com/flyteorg/flytectl/pkg/docker"

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
	containerFlyteSource                  = "/flyteorg/share/kustomization.yaml"
	generatedManifest                     = "/flyteorg/share/flyte_generated.yaml"
	flyteMinimumVersionSupported          = "v0.14.0"
	flyteMinimumVersionSupportedKustomize = "v0.15.1"
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
		isReady := docker.WaitForSandbox(reader, docker.SuccessMessage)
		if isReady {
			printExistingSandboxMessage()
			fmt.Println("")
			fmt.Printf("Please visit https://github.com/flyteorg/flytesnacks for more example %v \n", emoji.Rocket)
			fmt.Printf("Register all flytesnacks example by running 'flytectl register examples  -d development  -p flytesnacks' \n")
			return nil
		}
	}
	return fmt.Errorf("sandbox start failed. Please check troubleshooting documentation https://docs.flyte.org/en/latest/community/troubleshoot.html#troubleshoot")
}

func startSandbox(ctx context.Context, cli docker.Docker, reader io.Reader) (*bufio.Scanner, error) {
	fmt.Printf("%v Bootstrapping a brand new flyte cluster... %v %v\n", emoji.FactoryWorker, emoji.Hammer, emoji.Wrench)

	if err := docker.RemoveSandbox(ctx, cli, reader); err != nil {
		if err.Error() != clierrors.ErrSandboxExists {
			return nil, err
		}

		fmt.Printf("Existing details of your sandbox:")
		printExistingSandboxMessage()
		return nil, nil
	}

	if err := util.SetupFlyteDir(); err != nil {
		return nil, err
	}

	templateValues := configutil.ConfigTemplateSpec{
		Host:     "localhost:30081",
		Insecure: true,
	}

	templateStr := configutil.GetSandboxTemplate()
	if err := configutil.SetupConfig(configutil.FlytectlSandboxConfig, templateStr, templateValues); err != nil {
		return nil, err
	}

	volumes := docker.Volumes
	if vol, err := mountVolume(sandboxConfig.DefaultConfig.Source, docker.Source); err != nil {
		return nil, err
	} else if vol != nil {
		volumes = append(volumes, *vol)
	}

	if len(sandboxConfig.DefaultConfig.Version) > 0 {
		isGreater, err := util.IsVersionGreaterThanEqual(sandboxConfig.DefaultConfig.Version, flyteMinimumVersionSupported)
		if err != nil {
			return nil, err
		}
		if !isGreater {
			logger.Infof(ctx, "version flag only supported after with flyte %s+ release", flyteMinimumVersionSupported)
		}

		if err := downloadFlyteManifest(sandboxConfig.DefaultConfig.Version); err != nil {
			return nil, err
		}

		if vol, err := mountVolume(FlyteManifest, generatedManifest); err != nil {
			return nil, err
		} else if vol != nil {
			volumes = append(volumes, *vol)
		}
	}

	if len(sandboxConfig.DefaultConfig.Kustomize) > 0 {
		version := sandboxConfig.DefaultConfig.Version
		if len(sandboxConfig.DefaultConfig.Version) == 0 {
			release, err := util.GetLatestVersion("flyte")
			if err != nil {
				return nil, err
			}
			version = release.GetTagName()
		}
		isGreater, err := util.IsVersionGreaterThanEqual(version, flyteMinimumVersionSupportedKustomize)
		if err != nil {
			return nil, err
		}
		if isGreater {
			if vol, err := mountVolume(sandboxConfig.DefaultConfig.Kustomize, containerFlyteSource); err != nil {
				return nil, err
			} else if vol != nil {
				volumes = append(volumes, *vol)
			}
		} else {
			logger.Infof(ctx, "kustomize flag only supported after with flyte %s release", flyteMinimumVersionSupportedKustomize)
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

func printExistingSandboxMessage() {
	kubeconfig := strings.Join([]string{
		"$KUBECONFIG",
		f.FilePathJoin(f.UserHomeDir(), ".kube", "config"),
		docker.Kubeconfig,
	}, ":")
	fmt.Printf("%v %v %v %v %v \n", emoji.ManTechnologist, docker.SuccessMessage, emoji.Rocket, emoji.Rocket, emoji.PartyPopper)
	fmt.Printf("Add KUBECONFIG and FLYTECTL_CONFIG to your environment variable \n")
	fmt.Printf("export KUBECONFIG=%v \n", kubeconfig)
	fmt.Printf("export FLYTECTL_CONFIG=%v \n", configutil.FlytectlConfig)
}
