package docker

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"

	//"github.com/docker/go-connections/nat"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/flyteorg/flytectl/pkg/docker/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/docker/docker/api/types"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	u "github.com/flyteorg/flytectl/cmd/testutils"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"

	"github.com/stretchr/testify/assert"
)

var (
	cmdCtx     cmdCore.CommandContext
	containers []types.Container
)

func setupSandbox() {
	mockAdminClient := u.MockClient
	cmdCtx = cmdCore.NewCommandContext(mockAdminClient, u.MockOutStream)
	_ = SetupFlyteDir()
	container1 := types.Container{
		ID: "FlyteSandboxClusterName",
		Names: []string{
			FlyteSandboxClusterName,
		},
	}
	containers = append(containers, container1)
}

func TestConfigCleanup(t *testing.T) {
	_, err := os.Stat(f.FilePathJoin(f.UserHomeDir(), ".flyte"))
	if os.IsNotExist(err) {
		_ = os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte"), 0755)
	}
	_ = ioutil.WriteFile(FlytectlConfig, []byte("string"), 0600)
	_ = ioutil.WriteFile(Kubeconfig, []byte("string"), 0600)

	err = ConfigCleanup()
	assert.Nil(t, err)

	_, err = os.Stat(FlytectlConfig)
	check := os.IsNotExist(err)
	assert.Equal(t, check, true)

	_, err = os.Stat(Kubeconfig)
	check = os.IsNotExist(err)
	assert.Equal(t, check, true)
	_ = ConfigCleanup()
}

func TestSetupFlytectlConfig(t *testing.T) {
	_, err := os.Stat(f.FilePathJoin(f.UserHomeDir(), ".flyte"))
	if os.IsNotExist(err) {
		_ = os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte"), 0755)
	}
	err = SetupFlyteDir()
	assert.Nil(t, err)
	err = GetFlyteSandboxConfig()
	assert.Nil(t, err)
	_, err = os.Stat(FlytectlConfig)
	assert.Nil(t, err)
	check := os.IsNotExist(err)
	assert.Equal(t, check, false)
	_ = ConfigCleanup()
}

func TestGetSandbox(t *testing.T) {
	setupSandbox()
	mockDocker := &mocks.Docker{}
	context := context.Background()

	// Verify the attributes
	mockDocker.OnContainerList(context, types.ContainerListOptions{All: true}).Return(containers, nil)
	c := GetSandbox(context, mockDocker)
	assert.Equal(t, c.Names[0], FlyteSandboxClusterName)
}

func TestRemoveSandbox(t *testing.T) {
	setupSandbox()
	mockDocker := &mocks.Docker{}
	context := context.Background()

	// Verify the attributes
	mockDocker.OnContainerList(context, types.ContainerListOptions{All: true}).Return(containers, nil)
	mockDocker.OnContainerRemove(context, mock.Anything, types.ContainerRemoveOptions{Force: true}).Return(nil)
	err := RemoveSandbox(context, mockDocker, strings.NewReader("y"))
	assert.Nil(t, err)
}

func TestRemoveSandboxWithNo(t *testing.T) {
	setupSandbox()
	mockDocker := &mocks.Docker{}
	context := context.Background()

	// Verify the attributes
	mockDocker.OnContainerList(context, types.ContainerListOptions{All: true}).Return(containers, nil)
	mockDocker.OnContainerRemove(context, mock.Anything, types.ContainerRemoveOptions{Force: true}).Return(nil)
	err := RemoveSandbox(context, mockDocker, strings.NewReader("n"))
	assert.Nil(t, err)
}

func TestPullDockerImage(t *testing.T) {
	t.Run("Successful pull", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		context := context.Background()
		// Verify the attributes
		mockDocker.OnImagePullMatch(context, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		err := PullDockerImage(context, mockDocker, "nginx")
		assert.Nil(t, err)
	})

	t.Run("Error in pull", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		context := context.Background()
		// Verify the attributes
		mockDocker.OnImagePullMatch(context, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, fmt.Errorf("error"))
		err := PullDockerImage(context, mockDocker, "nginx")
		assert.NotNil(t, err)
	})

}

func TestStartContainer(t *testing.T) {
	p1, p2, _ := GetSandboxPorts()

	volumes := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: f.FilePathJoin(f.UserHomeDir(), ".flyte"),
			Target: K3sDir,
		},
	}

	t.Run("Success", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		context := context.Background()

		// Verify the attributes
		mockDocker.OnContainerCreate(context, &container.Config{
			Env:          Environment,
			Image:        ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(context, "Hello", types.ContainerStartOptions{}).Return(nil)
		id, err := StartContainer(context, mockDocker, volumes, p1, p2, "nginx", ImageName)
		assert.Nil(t, err)
		assert.Greater(t, len(id), 0)
		assert.Equal(t, id, "Hello")
	})

	t.Run("Error in create", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		context := context.Background()

		// Verify the attributes
		mockDocker.OnContainerCreate(context, &container.Config{
			Env:          Environment,
			Image:        ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "",
		}, fmt.Errorf("error"))
		mockDocker.OnContainerStart(context, "Hello", types.ContainerStartOptions{}).Return(nil)
		id, err := StartContainer(context, mockDocker, volumes, p1, p2, "nginx", ImageName)
		assert.NotNil(t, err)
		assert.Equal(t, len(id), 0)
		assert.Equal(t, id, "")
	})

	t.Run("Error in start", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		context := context.Background()

		// Verify the attributes
		mockDocker.OnContainerCreate(context, &container.Config{
			Env:          Environment,
			Image:        ImageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(context, "Hello", types.ContainerStartOptions{}).Return(fmt.Errorf("error"))
		id, err := StartContainer(context, mockDocker, volumes, p1, p2, "nginx", ImageName)
		assert.NotNil(t, err)
		assert.Equal(t, len(id), 0)
		assert.Equal(t, id, "")
	})
}

func TestWatchError(t *testing.T) {
	setupSandbox()
	mockDocker := &mocks.Docker{}
	context := context.Background()
	errCh := make(chan error)
	bodyStatus := make(chan container.ContainerWaitOKBody)
	mockDocker.OnContainerWaitMatch(context, mock.Anything, container.WaitConditionNotRunning).Return(bodyStatus, errCh)
	_, err := WatchError(context, mockDocker, "test")
	assert.NotNil(t, err)
}

func TestReadLogs(t *testing.T) {
	setupSandbox()

	t.Run("Success", func(t *testing.T) {
		mockDocker := &mocks.Docker{}
		context := context.Background()
		mockDocker.OnContainerLogsMatch(context, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		_, err := ReadLogs(context, mockDocker, "test")
		assert.Nil(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		mockDocker := &mocks.Docker{}
		context := context.Background()
		mockDocker.OnContainerLogsMatch(context, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, fmt.Errorf("error"))
		_, err := ReadLogs(context, mockDocker, "test")
		assert.NotNil(t, err)
	})
}

func TestWaitForSandbox(t *testing.T) {
	setupSandbox()
	reader := bufio.NewScanner(strings.NewReader("hello"))
	check := WaitForSandbox(reader, "hello")
	assert.Equal(t, check, true)

}
