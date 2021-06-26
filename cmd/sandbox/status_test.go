package sandbox

import (
	"context"
	"io"
	"testing"

	"github.com/docker/docker/api/types/filters"

	"github.com/docker/docker/api/types"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/docker"
	"github.com/flyteorg/flytectl/pkg/docker/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSandboxStatus(t *testing.T) {
	f := filters.NewArgs()
	f.Add("ancestor", docker.ImageName)
	t.Run("Sandbox status with zero result", func(t *testing.T) {
		ctx := context.Background()
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
		mockDocker := &mocks.Docker{}
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true, Filters: f}).Return([]types.Container{}, nil)
		docker.Client = mockDocker
		err := sandboxClusterStatus(ctx, []string{}, cmdCtx)
		assert.Nil(t, err)
	})
	t.Run("Sandbox status with running sandbox", func(t *testing.T) {
		ctx := context.Background()
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
		mockDocker := &mocks.Docker{}
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true, Filters: f}).Return([]types.Container{
			{
				ID: docker.FlyteSandboxClusterName,
				Names: []string{
					docker.FlyteSandboxClusterName,
				},
			},
		}, nil)
		docker.Client = mockDocker
		err := sandboxClusterStatus(ctx, []string{}, cmdCtx)
		assert.Nil(t, err)
	})
	t.Run("Sandbox status with name", func(t *testing.T) {
		ctx := context.Background()
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
		mockDocker := &mocks.Docker{}
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{
			{
				ID:    docker.FlyteSandboxClusterName,
				Image: "cr.flyte.org/flyteorg/flyte-sandbox:dind",
				Names: []string{
					docker.FlyteSandboxClusterName,
				},
				Ports: []types.Port{
					{PrivatePort: 30086, PublicPort: 30086},
					{PrivatePort: 30084, PublicPort: 30084},
					{PrivatePort: 30081, PublicPort: 30081},
					{PrivatePort: 30082, PublicPort: 30082},
				},
			},
		}, nil)
		docker.Client = mockDocker
		err := sandboxClusterStatus(ctx, []string{"flyte-sandbox"}, cmdCtx)
		assert.Nil(t, err)
	})
	t.Run("Sandbox status with a name that is not available", func(t *testing.T) {
		ctx := context.Background()
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
		mockDocker := &mocks.Docker{}
		mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{
			{
				ID:    docker.FlyteSandboxClusterName,
				Image: "cr.flyte.org/flyteorg/flyte-sandbox:dind",
				Names: []string{
					docker.FlyteSandboxClusterName,
				},
				Ports: []types.Port{
					{PrivatePort: 30086, PublicPort: 30086},
					{PrivatePort: 30084, PublicPort: 30084},
					{PrivatePort: 30081, PublicPort: 30081},
				},
			},
		}, nil)
		docker.Client = mockDocker
		err := sandboxClusterStatus(ctx, []string{"test"}, cmdCtx)
		assert.Nil(t, err)
	})
}
