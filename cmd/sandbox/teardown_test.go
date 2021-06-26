package sandbox

import (
	"context"
	"fmt"
	"io"
	"testing"

	sandboxConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/sandbox"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"

	"github.com/docker/docker/api/types"
	"github.com/flyteorg/flytectl/pkg/docker"
	"github.com/flyteorg/flytectl/pkg/docker/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTearDownFunc(t *testing.T) {

	var containers []types.Container
	container1 := types.Container{
		ID: "FlyteSandboxClusterName",
		Names: []string{
			fmt.Sprintf("/%v", docker.FlyteSandboxClusterName),
		},
	}
	containers = append(containers, container1)

	t.Run("Success on teardown", func(t *testing.T) {
		ctx := context.Background()
		mockDocker := &mocks.Docker{}
		mockDocker.OnContainerListMatch(ctx, types.ContainerListOptions{All: true}).Return(containers, nil)
		mockDocker.OnContainerRemoveMatch(ctx, mock.Anything, types.ContainerRemoveOptions{Force: true}).Return(nil)
		docker.Client = mockDocker
		err := tearDownSandbox(ctx, mockDocker)
		assert.Nil(t, err)
	})

	t.Run("Success on teardown sandbox with name", func(t *testing.T) {
		ctx := context.Background()
		sandboxConfig.DefaultConfig.Name = "test"
		mockDocker := &mocks.Docker{}
		mockDocker.OnContainerListMatch(ctx, types.ContainerListOptions{All: true}).Return(containers, nil)
		mockDocker.OnContainerRemoveMatch(ctx, mock.Anything, types.ContainerRemoveOptions{Force: true}).Return(fmt.Errorf("error"))
		docker.Client = mockDocker
		err := Execute(ctx, mockDocker, []string{"ls -al"})
		assert.Nil(t, err)
	})
	t.Run("Success on teardown command", func(t *testing.T) {
		mockOutStream := new(io.Writer)
		ctx := context.Background()
		cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
		mockDocker := &mocks.Docker{}
		mockDocker.OnContainerListMatch(ctx, types.ContainerListOptions{All: true}).Return(containers, nil)
		mockDocker.OnContainerRemoveMatch(ctx, mock.Anything, types.ContainerRemoveOptions{Force: true}).Return(fmt.Errorf("error"))
		docker.Client = mockDocker
		err := teardownSandboxCluster(ctx, []string{}, cmdCtx)
		assert.Nil(t, err)
	})

}
