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
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/enescakir/emoji"
	"github.com/flyteorg/flytectl/clierrors"
	sandboxCmdConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/sandbox"
	"github.com/flyteorg/flytectl/pkg/configutil"
	"github.com/flyteorg/flytectl/pkg/docker"
	"github.com/flyteorg/flytectl/pkg/github"
	"github.com/flyteorg/flytectl/pkg/k8s"
	"github.com/flyteorg/flytectl/pkg/util"
	"github.com/kataras/tablewriter"
	corev1api "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	flyteNamespace       = "flyte"
	diskPressureTaint    = "node.kubernetes.io/disk-pressure"
	taintEffect          = "NoSchedule"
	sandboxContextName   = "flyte-sandbox"
	sandboxDockerContext = "default"
	K8sEndpoint          = "https://127.0.0.1:6443"
	sandboxK8sEndpoint   = "https://127.0.0.1:30086"
	sandboxImageName     = "cr.flyte.org/flyteorg/flyte-sandbox"
	demoImageName        = "cr.flyte.org/flyteorg/flyte-sandbox-ultra"
)

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

func WatchFlyteDeployment(ctx context.Context, appsClient corev1.CoreV1Interface) error {
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
		} else {
			table.Append([]string{"k8s: This might take a little bit", "Bootstrapping", ""})
			table.Render()
		}

		time.Sleep(40 * time.Second)
	}

	return nil
}

func MountVolume(file, destination string) (*mount.Mount, error) {
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

func UpdateLocalKubeContext(dockerCtx string, contextName string, kubeConfigPath string) error {
	srcConfigAccess := &clientcmd.PathOptions{
		GlobalFile:   kubeConfigPath,
		LoadingRules: clientcmd.NewDefaultClientConfigLoadingRules(),
	}
	k8sCtxMgr := k8s.NewK8sContextManager()
	return k8sCtxMgr.CopyContext(srcConfigAccess, dockerCtx, contextName)
}

func startSandbox(ctx context.Context, cli docker.Docker, g github.GHRepoService, reader io.Reader, sandboxConfig *sandboxCmdConfig.Config, defaultImageName string, defaultImagePrefix string, exposedPorts map[nat.Port]struct{}, portBindings map[nat.Port][]nat.PortBinding, consolePort int) (*bufio.Scanner, error) {
	fmt.Printf("%v Bootstrapping a brand new flyte cluster... %v %v\n", emoji.FactoryWorker, emoji.Hammer, emoji.Wrench)

	if err := docker.RemoveSandbox(ctx, cli, reader); err != nil {
		if err.Error() != clierrors.ErrSandboxExists {
			return nil, err
		}
		fmt.Printf("Existing details of your sandbox")
		util.PrintSandboxMessage(consolePort)
		return nil, nil
	}

	templateValues := configutil.ConfigTemplateSpec{
		Host:     "localhost:30080",
		Insecure: true,
		Console:  fmt.Sprintf("http://localhost:%d", consolePort),
	}
	if err := configutil.SetupConfig(configutil.FlytectlConfig, configutil.GetTemplate(), templateValues); err != nil {
		return nil, err
	}

	volumes := docker.Volumes
	// Mount this even though it should no longer be necessary. This is for user code
	if vol, err := MountVolume(sandboxConfig.DeprecatedSource, docker.Source); err != nil {
		return nil, err
	} else if vol != nil {
		volumes = append(volumes, *vol)
	}

	// This is the state directory mount, where flyte sandbox will write postgres/blobs, configs
	if vol, err := MountVolume(docker.FlyteStateDir, docker.StateDirMountDest); err != nil {
		return nil, err
	} else if vol != nil {
		volumes = append(volumes, *vol)
	}

	sandboxImage := sandboxConfig.Image
	if len(sandboxImage) == 0 {
		image, version, err := github.GetFullyQualifiedImageName(defaultImagePrefix, sandboxConfig.Version, defaultImageName, sandboxConfig.Prerelease, g)
		if err != nil {
			return nil, err
		}
		sandboxImage = image
		fmt.Printf("%s Fully Qualified image\n", image)
		fmt.Printf("%v Running Flyte %s release\n", emoji.Whale, version)
	}
	fmt.Printf("%v pulling docker image for release %s\n", emoji.Whale, sandboxImage)
	if err := docker.PullDockerImage(ctx, cli, sandboxImage, sandboxConfig.ImagePullPolicy, sandboxConfig.ImagePullOptions); err != nil {
		return nil, err
	}
	sandboxEnv := sandboxConfig.Env
	if sandboxConfig.Dev {
		sandboxEnv = append(sandboxEnv, "FLYTE_DEV=True")
	}

	fmt.Printf("%v booting Flyte-sandbox container\n", emoji.FactoryWorker)
	ID, err := docker.StartContainer(ctx, cli, volumes, exposedPorts, portBindings, docker.FlyteSandboxClusterName,
		sandboxImage, sandboxEnv)

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

func primeFlytekitPod(ctx context.Context, podService corev1.PodInterface) {
	_, err := podService.Create(ctx, &corev1api.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "py39-cacher",
		},
		Spec: corev1api.PodSpec{
			RestartPolicy: corev1api.RestartPolicyNever,
			Containers: []corev1api.Container{
				{

					Name:    "flytekit",
					Image:   "ghcr.io/flyteorg/flytekit:py3.9-latest",
					Command: []string{"echo"},
					Args:    []string{"Flyte"},
				},
			},
		},
	}, v1.CreateOptions{})
	if err != nil {
		fmt.Printf("Failed to create primer pod - %s", err)
	}
}

func StartCluster(ctx context.Context, args []string, sandboxConfig *sandboxCmdConfig.Config, primePod bool, defaultImageName string, defaultImagePrefix string, exposedPorts map[nat.Port]struct{}, portBindings map[nat.Port][]nat.PortBinding, consolePort int) error {
	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	ghRepo := github.GetGHRepoService()
	if err := util.CreatePathAndFile(docker.Kubeconfig); err != nil {
		return err
	}

	reader, err := startSandbox(ctx, cli, ghRepo, os.Stdin, sandboxConfig, defaultImageName, defaultImagePrefix, exposedPorts, portBindings, consolePort)
	if err != nil {
		return err
	}

	if reader != nil {
		var k8sClient k8s.K8s
		err = retry.Do(
			func() error {
				// This should wait for the kubeconfig file being there.
				k8sClient, err = k8s.GetK8sClient(docker.Kubeconfig, K8sEndpoint)
				return err
			},
			retry.Attempts(10),
		)
		if err != nil {
			return err
		}

		// This will copy the kubeconfig from where k3s writes it () to the main file.
		if err = UpdateLocalKubeContext(sandboxDockerContext, sandboxContextName, docker.Kubeconfig); err != nil {
			return err
		}

		// Live-ness check
		err = retry.Do(
			func() error {
				// Have to get a new client every time because you run into x509 errors if not
				fmt.Println("Waiting for cluster to come up...")
				k8sClient, err = k8s.GetK8sClient(docker.Kubeconfig, K8sEndpoint)
				req := k8sClient.CoreV1().RESTClient().Get()
				req = req.RequestURI("livez")
				res := req.Do(ctx)
				return res.Error()
			},
			retry.Attempts(15),
		)
		if err != nil {
			return err
		}

		// Readiness check
		err = retry.Do(
			func() error {
				// No need to refresh client here
				req := k8sClient.CoreV1().RESTClient().Get()
				req = req.RequestURI("readyz")
				res := req.Do(ctx)
				return res.Error()
			},
			retry.Attempts(10),
		)
		if err != nil {
			return err
		}

		// Watch for Flyte Deployment
		if err := WatchFlyteDeployment(ctx, k8sClient.CoreV1()); err != nil {
			return err
		}
		if primePod {
			primeFlytekitPod(ctx, k8sClient.CoreV1().Pods("default"))
		}
	}
	return nil
}

// StartClusterForSandbox is the code for the original multi deploy version of sandbox, should be removed once we
// document the new development experience for plugins.
func StartClusterForSandbox(ctx context.Context, args []string, sandboxConfig *sandboxCmdConfig.Config, primePod bool, defaultImageName string, defaultImagePrefix string, exposedPorts map[nat.Port]struct{}, portBindings map[nat.Port][]nat.PortBinding, consolePort int) error {
	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}

	ghRepo := github.GetGHRepoService()

	if err := util.CreatePathAndFile(docker.SandboxKubeconfig); err != nil {
		return err
	}

	reader, err := startSandbox(ctx, cli, ghRepo, os.Stdin, sandboxConfig, defaultImageName, defaultImagePrefix, exposedPorts, portBindings, consolePort)
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
				k8sClient, err = k8s.GetK8sClient(docker.SandboxKubeconfig, sandboxK8sEndpoint)
				return err
			},
			retry.Attempts(10),
		)
		if err != nil {
			return err
		}
		if err = UpdateLocalKubeContext(sandboxDockerContext, sandboxContextName, docker.SandboxKubeconfig); err != nil {
			return err
		}

		// TODO: This doesn't appear to correctly watch for the Flyte deployment but doesn't do so on master either.
		if err := WatchFlyteDeployment(ctx, k8sClient.CoreV1()); err != nil {
			return err
		}
		if primePod {
			primeFlytekitPod(ctx, k8sClient.CoreV1().Pods("default"))
		}

	}
	return nil
}

func DemoClusterInit(ctx context.Context, args []string, sandboxConfig *sandboxCmdConfig.Config) error {
	sandboxImagePrefix := "sha"

	// TODO: Add check and warning if the file already exists
	// TODO: Make sure the state folder is created

	cli, err := docker.GetDockerClient()
	if err != nil {
		return err
	}
	ghRepo := github.GetGHRepoService()

	// Determine and pull the image
	sandboxImage := sandboxConfig.Image
	if len(sandboxImage) == 0 {
		image, _, err := github.GetFullyQualifiedImageName(sandboxImagePrefix, sandboxConfig.Version, demoImageName, sandboxConfig.Prerelease, ghRepo)
		if err != nil {
			return err
		}
		sandboxImage = image
	}
	fmt.Printf("%v Fetching image %s\n", emoji.Whale, sandboxImage)
	err = docker.PullDockerImage(ctx, cli, sandboxImage, docker.ImagePullPolicyIfNotPresent, docker.ImagePullOptions{})
	if err != nil {
		return err
	}

	err = docker.CopyContainerFile(ctx, cli, "/opt/flyte/defaults.flyte.yaml", docker.FlyteBinaryConfig, "demo-init", sandboxImage)
	if err != nil {
		return err
	}
	return nil
}

func StartDemoCluster(ctx context.Context, args []string, sandboxConfig *sandboxCmdConfig.Config) error {
	primePod := true
	sandboxImagePrefix := "sha"
	exposedPorts, portBindings, err := docker.GetDemoPorts()
	if err != nil {
		return err
	}
	// TODO: Bring back dev later
	// K3s will automatically write the file specified by this var, which is mounted from user's local state dir.
	sandboxConfig.Env = append(sandboxConfig.Env, "K3S_KUBECONFIG_OUTPUT=/srv/flyte/kubeconfig")
	err = StartCluster(ctx, args, sandboxConfig, primePod, demoImageName, sandboxImagePrefix, exposedPorts, portBindings, util.DemoConsolePort)
	if err != nil {
		return err
	}
	util.PrintSandboxMessage(util.DemoConsolePort)
	return nil
}

func StartSandboxCluster(ctx context.Context, args []string, sandboxConfig *sandboxCmdConfig.Config) error {
	primePod := false
	demoImagePrefix := "dind"
	exposedPorts, portBindings, err := docker.GetSandboxPorts()
	if err != nil {
		return err
	}
	err = StartClusterForSandbox(ctx, args, sandboxConfig, primePod, sandboxImageName, demoImagePrefix, exposedPorts, portBindings, util.SandBoxConsolePort)
	if err != nil {
		return err
	}
	util.PrintSandboxMessage(util.SandBoxConsolePort)
	return nil
}
