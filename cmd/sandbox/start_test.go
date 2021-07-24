package sandbox

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	sandboxConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/sandbox"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/docker"
	"github.com/flyteorg/flytectl/pkg/docker/mocks"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	"github.com/flyteorg/flytectl/pkg/k8s"
	"github.com/flyteorg/flytectl/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

var content = `
apiVersion: v1
clusters:
- cluster:
    server: https://localhost:8080
    extensions:
    - name: client.authentication.k8s.io/exec
      extension:
        audience: foo
        other: bar
  name: foo-cluster
contexts:
- context:
    cluster: foo-cluster
    user: foo-user
    namespace: bar
  name: foo-context
current-context: foo-context
kind: Config
users:
- name: foo-user
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - arg-1
      - arg-2
      command: foo-command
      provideClusterInfo: true
`

var fakeNodeTaint = &corev1.Node{
	Status: corev1.NodeStatus{
		Conditions: []corev1.NodeCondition{
			{
				Type:   "MemoryPressure",
				Status: corev1.ConditionFalse,
				Reason: "KubeletHasSufficientMemory",
			},
			{
				Type:   "DiskPressure",
				Status: corev1.ConditionFalse,
				Reason: "KubeletHasDiskPressure",
			},
		},
	},
}

func TestStartSandboxFunc(t *testing.T) {
	p1, p2, _ := docker.GetSandboxPorts()
	client := testclient.NewSimpleClientset()
	k8s.Client = client
	assert.Nil(t, util.SetupFlyteDir())
	assert.Nil(t, os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte", "k3s"), os.ModePerm))
	assert.Nil(t, ioutil.WriteFile(docker.Kubeconfig, []byte(content), os.ModePerm))

	fakeDeployment := appsv1.Deployment{
		Status: appsv1.DeploymentStatus{
			AvailableReplicas: 0,
		},
	}
	fakeDeployment.SetName("flyte")
	fakeDeployment.SetName("flyte")

	t.Run("Successfully run sandbox cluster", func(t *testing.T) {
		ctx := context.Background()
		mockDocker := &mocks.Docker{}
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       docker.Volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		_, err := startSandbox(ctx, mockDocker, os.Stdin)
		assert.Nil(t, err)
	})
	t.Run("Successfully exit when sandbox cluster exist", func(t *testing.T) {
		ctx := context.Background()
		mockDocker := &mocks.Docker{}
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       docker.Volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{
			{
				ID: docker.FlyteSandboxClusterName,
				Names: []string{
					docker.FlyteSandboxClusterName,
				},
			},
		}, nil)
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		reader, err := startSandbox(ctx, mockDocker, strings.NewReader("n"))
		assert.Nil(t, err)
		assert.Nil(t, reader)
	})
	t.Run("Successfully run sandbox cluster with source code", func(t *testing.T) {
		ctx := context.Background()
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker := &mocks.Docker{}
		sandboxConfig.DefaultConfig.Source = f.UserHomeDir()
		volumes := docker.Volumes
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: sandboxConfig.DefaultConfig.Source,
			Target: docker.Source,
		})
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		_, err := startSandbox(ctx, mockDocker, os.Stdin)
		assert.Nil(t, err)
	})
	t.Run("Successfully run sandbox cluster with abs path of source code", func(t *testing.T) {
		ctx := context.Background()
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker := &mocks.Docker{}
		sandboxConfig.DefaultConfig.Source = "../"
		absPath, err := filepath.Abs(sandboxConfig.DefaultConfig.Source)
		assert.Nil(t, err)
		volumes := docker.Volumes
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: absPath,
			Target: docker.Source,
		})
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		_, err = startSandbox(ctx, mockDocker, os.Stdin)
		assert.Nil(t, err)
	})
	t.Run("Successfully run sandbox cluster with specific version", func(t *testing.T) {
		ctx := context.Background()
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker := &mocks.Docker{}
		sandboxConfig.DefaultConfig.Version = "v0.15.0"
		sandboxConfig.DefaultConfig.Source = ""
		volumes := docker.Volumes
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: FlyteManifest,
			Target: GeneratedManifest,
		})
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		_, err := startSandbox(ctx, mockDocker, os.Stdin)
		assert.Nil(t, err)
	})
	t.Run("Failed run sandbox cluster with wrong version", func(t *testing.T) {
		ctx := context.Background()
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker := &mocks.Docker{}
		sandboxConfig.DefaultConfig.Version = "v0.13.0"
		sandboxConfig.DefaultConfig.Source = ""
		volumes := docker.Volumes
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: FlyteManifest,
			Target: GeneratedManifest,
		})
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		_, err := startSandbox(ctx, mockDocker, os.Stdin)
		assert.NotNil(t, err)
	})
	t.Run("Error in pulling image", func(t *testing.T) {
		ctx := context.Background()
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker := &mocks.Docker{}
		sandboxConfig.DefaultConfig.Source = f.UserHomeDir()
		volumes := docker.Volumes
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: sandboxConfig.DefaultConfig.Source,
			Target: docker.Source,
		})
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, fmt.Errorf("error"))
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		_, err := startSandbox(ctx, mockDocker, os.Stdin)
		assert.NotNil(t, err)
	})
	t.Run("Error in  removing existing cluster", func(t *testing.T) {
		ctx := context.Background()
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker := &mocks.Docker{}
		sandboxConfig.DefaultConfig.Source = f.UserHomeDir()
		volumes := docker.Volumes
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: sandboxConfig.DefaultConfig.Source,
			Target: docker.Source,
		})
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{
			{
				ID: docker.FlyteSandboxClusterName,
				Names: []string{
					docker.FlyteSandboxClusterName,
				},
			},
		}, nil)
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		mockDocker.OnContainerRemove(ctx, mock.Anything, types.ContainerRemoveOptions{Force: true}).Return(fmt.Errorf("error"))
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		_, err := startSandbox(ctx, mockDocker, strings.NewReader("y"))
		assert.NotNil(t, err)
	})
	t.Run("Error in start container", func(t *testing.T) {
		ctx := context.Background()
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker := &mocks.Docker{}
		sandboxConfig.DefaultConfig.Source = ""
		sandboxConfig.DefaultConfig.Version = ""
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       docker.Volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, fmt.Errorf("error"))
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(fmt.Errorf("error"))
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		_, err := startSandbox(ctx, mockDocker, os.Stdin)
		assert.NotNil(t, err)
	})
	t.Run("Failed manifest", func(t *testing.T) {
		err := downloadFlyteManifest("v100.9.9")
		assert.NotNil(t, err)
	})
	t.Run("Error in reading logs", func(t *testing.T) {
		ctx := context.Background()
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker := &mocks.Docker{}
		sandboxConfig.DefaultConfig.Source = f.UserHomeDir()
		volumes := docker.Volumes
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: sandboxConfig.DefaultConfig.Source,
			Target: docker.Source,
		})
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, fmt.Errorf("error"))
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		_, err := startSandbox(ctx, mockDocker, os.Stdin)
		assert.NotNil(t, err)
	})
	t.Run("Error in list container", func(t *testing.T) {
		ctx := context.Background()
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker := &mocks.Docker{}
		sandboxConfig.DefaultConfig.Source = f.UserHomeDir()
		volumes := docker.Volumes
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: sandboxConfig.DefaultConfig.Source,
			Target: docker.Source,
		})
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, fmt.Errorf("error"))
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		_, err := startSandbox(ctx, mockDocker, os.Stdin)
		assert.Nil(t, err)
	})
	t.Run("Successfully run sandbox cluster command", func(t *testing.T) {
		mockOutStream := new(io.Writer)
		ctx := context.Background()
		cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
		mockDocker := &mocks.Docker{}
		errCh := make(chan error)
		_, err := client.AppsV1().Deployments("flyte").Create(ctx, &fakeDeployment, v1.CreateOptions{})
		if err != nil {
			t.Error(err)
		}
		fakeNodeTaint.SetName("flyte")
		_, err = client.CoreV1().Nodes().Create(ctx, fakeNodeTaint, v1.CreateOptions{})
		if err != nil {
			t.Error(err)
		}
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       docker.Volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		stringReader := strings.NewReader(docker.SuccessMessage)
		reader := ioutil.NopCloser(stringReader)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(reader, nil)
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		docker.Client = mockDocker
		sandboxConfig.DefaultConfig.Source = ""
		go func() {
			dep, err := client.AppsV1().Deployments("flyte").Get(ctx, "flyte", v1.GetOptions{})
			if err != nil {
				t.Error(err)
			}
			dep.Status.AvailableReplicas = 1
			time.Sleep(2 * time.Second)
			_, err = client.AppsV1().Deployments("flyte").Update(ctx, dep, v1.UpdateOptions{})
			if err != nil {
				t.Error(err)
			}
		}()
		err = startSandboxCluster(ctx, []string{}, cmdCtx)
		assert.NotNil(t, err)
	})
	t.Run("Error in running sandbox cluster command", func(t *testing.T) {
		mockOutStream := new(io.Writer)
		ctx := context.Background()
		cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
		mockDocker := &mocks.Docker{}
		errCh := make(chan error)
		bodyStatus := make(chan container.ContainerWaitOKBody)
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          docker.Environment,
			Image:        docker.ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       docker.Volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(fmt.Errorf("error"))
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, fmt.Errorf("error"))
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		stringReader := strings.NewReader(docker.SuccessMessage)
		reader := ioutil.NopCloser(stringReader)
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(reader, nil)
		mockDocker.OnContainerWaitMatch(ctx, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
		docker.Client = mockDocker
		sandboxConfig.DefaultConfig.Source = ""
		err := startSandboxCluster(ctx, []string{}, cmdCtx)
		assert.NotNil(t, err)
	})
}

func TestMonitorFlyteDeployment(t *testing.T) {
	t.Run("Monitor k8s deployment", func(t *testing.T) {
		ctx := context.Background()
		client := testclient.NewSimpleClientset()
		fakeDeployment := &appsv1.Deployment{
			Status: appsv1.DeploymentStatus{
				AvailableReplicas: 1,
			},
		}
		fakeDeployment.SetName("flyte")
		fakeDeployment.SetNamespace("flyte")
		_, err := client.AppsV1().Deployments("flyte").Create(ctx, fakeDeployment, v1.CreateOptions{})
		if err != nil {
			t.Error(err)
		}
		total, _, _ := monitorFlyteDeployment(ctx, client.AppsV1(), client.CoreV1().Nodes())

		assert.Equal(t, int64(1), total)
	})
}
