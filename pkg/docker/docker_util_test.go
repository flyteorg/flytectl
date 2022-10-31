package docker

import (
	"archive/tar"
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"

	"github.com/docker/docker/api/types/container"
	"github.com/flyteorg/flytectl/pkg/docker/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

var (
	containers []types.Container
	imageName  = "cr.flyte.org/flyteorg/flyte-sandbox"
)

func setupSandbox() {
	err := os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte"), os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	container1 := types.Container{
		ID: "FlyteSandboxClusterName",
		Names: []string{
			FlyteSandboxClusterName,
		},
	}
	containers = append(containers, container1)
}

func TestGetSandbox(t *testing.T) {
	setupSandbox()
	t.Run("Successfully get sandbox container", func(t *testing.T) {
		mockDocker := &mocks.Docker{}
		ctx := context.Background()

		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return(containers, nil)
		c, err := GetSandbox(ctx, mockDocker)
		assert.Equal(t, c.Names[0], FlyteSandboxClusterName)
		assert.Nil(t, err)
	})

	t.Run("Successfully get sandbox container with zero result", func(t *testing.T) {
		mockDocker := &mocks.Docker{}
		ctx := context.Background()

		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
		c, err := GetSandbox(ctx, mockDocker)
		assert.Nil(t, c)
		assert.Nil(t, err)
	})

	t.Run("Error in get sandbox container", func(t *testing.T) {
		mockDocker := &mocks.Docker{}
		ctx := context.Background()

		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return(containers, nil)
		mockDocker.OnContainerRemove(ctx, mock.Anything, types.ContainerRemoveOptions{Force: true}).Return(nil)
		err := RemoveSandbox(ctx, mockDocker, strings.NewReader("y"))
		assert.Nil(t, err)
	})

}

func TestRemoveSandboxWithNoReply(t *testing.T) {
	setupSandbox()
	t.Run("Successfully remove sandbox container", func(t *testing.T) {
		mockDocker := &mocks.Docker{}
		ctx := context.Background()

		// Verify the attributes
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return(containers, nil)
		mockDocker.OnContainerRemove(ctx, mock.Anything, types.ContainerRemoveOptions{Force: true}).Return(nil)
		err := RemoveSandbox(ctx, mockDocker, strings.NewReader("n"))
		assert.NotNil(t, err)
	})

	t.Run("Successfully remove sandbox container with zero sandbox containers are running", func(t *testing.T) {
		mockDocker := &mocks.Docker{}
		ctx := context.Background()

		// Verify the attributes
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{}, nil)
		mockDocker.OnContainerRemove(ctx, mock.Anything, types.ContainerRemoveOptions{Force: true}).Return(nil)
		err := RemoveSandbox(ctx, mockDocker, strings.NewReader("n"))
		assert.Nil(t, err)
	})

}

func TestPullDockerImage(t *testing.T) {
	t.Run("Successfully pull image Always", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		ctx := context.Background()
		// Verify the attributes
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		err := PullDockerImage(ctx, mockDocker, "nginx:latest", ImagePullPolicyAlways, ImagePullOptions{})
		assert.Nil(t, err)
	})

	t.Run("Error in pull image", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		ctx := context.Background()
		// Verify the attributes
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, fmt.Errorf("error"))
		err := PullDockerImage(ctx, mockDocker, "nginx:latest", ImagePullPolicyAlways, ImagePullOptions{})
		assert.NotNil(t, err)
	})

	t.Run("Successfully pull image IfNotPresent", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		ctx := context.Background()
		// Verify the attributes
		mockDocker.OnImagePullMatch(ctx, mock.Anything, types.ImagePullOptions{}).Return(os.Stdin, nil)
		mockDocker.OnImageListMatch(ctx, types.ImageListOptions{}).Return([]types.ImageSummary{}, nil)
		err := PullDockerImage(ctx, mockDocker, "nginx:latest", ImagePullPolicyIfNotPresent, ImagePullOptions{})
		assert.Nil(t, err)
	})

	t.Run("Successfully pull image Never", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		ctx := context.Background()
		err := PullDockerImage(ctx, mockDocker, "nginx:latest", ImagePullPolicyNever, ImagePullOptions{})
		assert.Nil(t, err)
	})
}

func TestStartContainer(t *testing.T) {
	p1, p2, _ := GetSandboxPorts()

	t.Run("Successfully create a container", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		ctx := context.Background()

		// Verify the attributes
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          Environment,
			Image:        imageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       Volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		id, err := StartContainer(ctx, mockDocker, Volumes, p1, p2, "nginx", imageName, nil)
		assert.Nil(t, err)
		assert.Greater(t, len(id), 0)
		assert.Equal(t, id, "Hello")
	})

	t.Run("Successfully create a container with Env", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		ctx := context.Background()
		// Setup additional env
		additionalEnv := []string{"a=1", "b=2"}
		expectedEnv := append(Environment, "a=1")
		expectedEnv = append(expectedEnv, "b=2")

		// Verify the attributes
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          expectedEnv,
			Image:        imageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       Volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		id, err := StartContainer(ctx, mockDocker, Volumes, p1, p2, "nginx", imageName, additionalEnv)
		assert.Nil(t, err)
		assert.Greater(t, len(id), 0)
		assert.Equal(t, id, "Hello")
		assert.Equal(t, expectedEnv, Environment)
	})

	t.Run("Error in creating container", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		ctx := context.Background()

		// Verify the attributes
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          Environment,
			Image:        imageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       Volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "",
		}, fmt.Errorf("error"))
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(nil)
		id, err := StartContainer(ctx, mockDocker, Volumes, p1, p2, "nginx", imageName, nil)
		assert.NotNil(t, err)
		assert.Equal(t, len(id), 0)
		assert.Equal(t, id, "")
	})

	t.Run("Error in start of a container", func(t *testing.T) {
		setupSandbox()
		mockDocker := &mocks.Docker{}
		ctx := context.Background()

		// Verify the attributes
		mockDocker.OnContainerCreate(ctx, &container.Config{
			Env:          Environment,
			Image:        imageName,
			Tty:          false,
			ExposedPorts: p1,
		}, &container.HostConfig{
			Mounts:       Volumes,
			PortBindings: p2,
			Privileged:   true,
		}, nil, nil, mock.Anything).Return(container.ContainerCreateCreatedBody{
			ID: "Hello",
		}, nil)
		mockDocker.OnContainerStart(ctx, "Hello", types.ContainerStartOptions{}).Return(fmt.Errorf("error"))
		id, err := StartContainer(ctx, mockDocker, Volumes, p1, p2, "nginx", imageName, nil)
		assert.NotNil(t, err)
		assert.Equal(t, len(id), 0)
		assert.Equal(t, id, "")
	})
}

func TestReadLogs(t *testing.T) {
	setupSandbox()

	t.Run("Successfully read logs", func(t *testing.T) {
		mockDocker := &mocks.Docker{}
		ctx := context.Background()
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, nil)
		_, err := ReadLogs(ctx, mockDocker, "test")
		assert.Nil(t, err)
	})

	t.Run("Error in reading logs", func(t *testing.T) {
		mockDocker := &mocks.Docker{}
		ctx := context.Background()
		mockDocker.OnContainerLogsMatch(ctx, mock.Anything, types.ContainerLogsOptions{
			ShowStderr: true,
			ShowStdout: true,
			Timestamps: true,
			Follow:     true,
		}).Return(nil, fmt.Errorf("error"))
		_, err := ReadLogs(ctx, mockDocker, "test")
		assert.NotNil(t, err)
	})
}

func TestWaitForSandbox(t *testing.T) {
	setupSandbox()
	t.Run("Successfully read logs ", func(t *testing.T) {
		reader := bufio.NewScanner(strings.NewReader("hello \n Flyte"))

		check := WaitForSandbox(reader, "Flyte")
		assert.Equal(t, true, check)
	})

	t.Run("Error in reading logs ", func(t *testing.T) {
		reader := bufio.NewScanner(strings.NewReader(""))
		check := WaitForSandbox(reader, "Flyte")
		assert.Equal(t, false, check)
	})
}

func TestDockerClient(t *testing.T) {
	t.Run("Successfully get docker mock client", func(t *testing.T) {
		mockDocker := &mocks.Docker{}
		Client = mockDocker
		cli, err := GetDockerClient()
		assert.Nil(t, err)
		assert.NotNil(t, cli)
	})
	t.Run("Successfully get docker client", func(t *testing.T) {
		Client = nil
		cli, err := GetDockerClient()
		assert.Nil(t, err)
		assert.NotNil(t, cli)
	})
}

func TestDockerExec(t *testing.T) {
	t.Run("Successfully exec command in container", func(t *testing.T) {
		ctx := context.Background()
		mockDocker := &mocks.Docker{}
		Client = mockDocker
		c := ExecConfig
		c.Cmd = []string{"ls"}
		mockDocker.OnContainerExecCreateMatch(ctx, mock.Anything, c).Return(types.IDResponse{}, nil)
		_, err := ExecCommend(ctx, mockDocker, "test", []string{"ls"})
		assert.Nil(t, err)
	})
	t.Run("Failed exec command in container", func(t *testing.T) {
		ctx := context.Background()
		mockDocker := &mocks.Docker{}
		Client = mockDocker
		c := ExecConfig
		c.Cmd = []string{"ls"}
		mockDocker.OnContainerExecCreateMatch(ctx, mock.Anything, c).Return(types.IDResponse{}, fmt.Errorf("test"))
		_, err := ExecCommend(ctx, mockDocker, "test", []string{"ls"})
		assert.NotNil(t, err)
	})
}

func TestInspectExecResp(t *testing.T) {
	t.Run("Failed exec command in container", func(t *testing.T) {
		ctx := context.Background()
		mockDocker := &mocks.Docker{}
		Client = mockDocker
		c := ExecConfig
		c.Cmd = []string{"ls"}
		reader := bufio.NewReader(strings.NewReader("test"))

		mockDocker.OnContainerExecInspectMatch(ctx, mock.Anything).Return(types.ContainerExecInspect{}, nil)
		mockDocker.OnContainerExecAttachMatch(ctx, mock.Anything, types.ExecStartCheck{}).Return(types.HijackedResponse{
			Reader: reader,
		}, fmt.Errorf("err"))

		err := InspectExecResp(ctx, mockDocker, "test")
		assert.NotNil(t, err)
	})
	t.Run("Successfully exec command in container", func(t *testing.T) {
		ctx := context.Background()
		mockDocker := &mocks.Docker{}
		Client = mockDocker
		c := ExecConfig
		c.Cmd = []string{"ls"}
		reader := bufio.NewReader(strings.NewReader("test"))

		mockDocker.OnContainerExecAttachMatch(ctx, mock.Anything, types.ExecStartCheck{}).Return(types.HijackedResponse{
			Reader: reader,
		}, nil)

		err := InspectExecResp(ctx, mockDocker, "test")
		assert.Nil(t, err)
	})

}

func TestDemoPorts(t *testing.T) {
	_, ports, _ := GetDemoPorts()
	assert.Equal(t, 5, len(ports))
}

func TestTaring(t *testing.T) {
	// Create a fake tar file in tmp.
	text := "a: b"
	fo, err := os.CreateTemp("", "sampledata")
	assert.NoError(t, err)
	tarWriter := tar.NewWriter(fo)
	err = tarWriter.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg,
		Name:     "flyte.yaml",
		Size:     4,
		Mode:     0640,
		ModTime:  time.Unix(1245206587, 0),
	})
	assert.NoError(t, err)
	cnt, err := tarWriter.Write([]byte(text))
	assert.NoError(t, err)
	assert.Equal(t, 4, cnt)
	tarWriter.Close()
	fo.Close()

	t.Run("Basic testing", func(t *testing.T) {
		destFile, err := os.CreateTemp("", "sampledata")
		assert.NoError(t, err)
		reader, err := os.Open(fo.Name())
		err = ExtractTar(reader, destFile.Name())
		assert.NoError(t, err)
		fileBytes, err := os.ReadFile(destFile.Name())
		assert.NoError(t, err)
		readString := string(fileBytes)
		assert.Equal(t, text, readString)

		// Try to extract the file we just extracted again. It's not a tar file obviously so it should error
		reader, err = os.Open(destFile.Name())
		err = ExtractTar(reader, destFile.Name())
		assert.Errorf(t, err, "unexpected EOF")
	})
}

func TestTarBadHeader(t *testing.T) {
	// Create a fake tar file in tmp.
	fo, err := os.CreateTemp("", "sampledata")
	assert.NoError(t, err)
	tarWriter := tar.NewWriter(fo)
	// Write a symlink, we should not know how to parse.
	err = tarWriter.WriteHeader(&tar.Header{
		Typeflag: tar.TypeLink,
		Name:     "flyte.yaml",
		Size:     4,
		Mode:     0640,
		ModTime:  time.Unix(1245206587, 0),
	})
	assert.NoError(t, err)
	tarWriter.Close()
	fo.Close()

	t.Run("Basic testing", func(t *testing.T) {
		destFile, err := os.CreateTemp("", "sampledata")
		assert.NoError(t, err)
		reader, err := os.Open(fo.Name())
		err = ExtractTar(reader, destFile.Name())
		assert.Errorf(t, err, "ExtractTarGz: unknown type")
	})
}
