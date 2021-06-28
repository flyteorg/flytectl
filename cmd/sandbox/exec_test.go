package sandbox

import (
	"bufio"
	"context"
	"io"
	"strings"
	"testing"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
	"github.com/flyteorg/flytectl/pkg/docker"
	"github.com/flyteorg/flytectl/pkg/docker/mocks"
	"github.com/stretchr/testify/mock"
)

func TestSandboxClusterExec(t *testing.T) {
	mockDocker := &mocks.Docker{}
	mockOutStream := new(io.Writer)
	ctx := context.Background()
	cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
	reader := bufio.NewReader(strings.NewReader("test"))

	mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{
		{
			ID: docker.FlyteSandboxClusterName,
			Names: []string{
				docker.FlyteSandboxClusterName,
			},
		},
	}, nil)
	docker.ExecConfig.Cmd = []string{"ls -al"}
	mockDocker.OnContainerExecCreateMatch(ctx, mock.Anything, docker.ExecConfig).Return(types.IDResponse{}, nil)
	mockDocker.OnContainerExecInspectMatch(ctx, mock.Anything).Return(types.ContainerExecInspect{}, nil)
	mockDocker.OnContainerExecAttachMatch(ctx, mock.Anything, types.ExecStartCheck{}).Return(types.HijackedResponse{
		Reader: reader,
	}, nil)
	docker.Client = mockDocker
	err := sandboxClusterExec(ctx, []string{"ls -al"}, cmdCtx)
	assert.Nil(t, err)
}

func TestSandboxClusterExecWithNoCommand(t *testing.T) {
	mockDocker := &mocks.Docker{}
	mockOutStream := new(io.Writer)
	ctx := context.Background()
	cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
	c := docker.ExecConfig
	c.Cmd = []string{"ls"}

	reader := bufio.NewReader(strings.NewReader("test"))

	mockDocker.OnContainerList(ctx, types.ContainerListOptions{All: true}).Return([]types.Container{
		{
			ID: docker.FlyteSandboxClusterName,
			Names: []string{
				docker.FlyteSandboxClusterName,
			},
		},
	}, nil)
	docker.ExecConfig.Cmd = []string{}
	mockDocker.OnContainerExecCreateMatch(ctx, mock.Anything, docker.ExecConfig).Return(types.IDResponse{}, nil)
	mockDocker.OnContainerExecInspectMatch(ctx, mock.Anything).Return(types.ContainerExecInspect{}, nil)
	mockDocker.OnContainerExecAttachMatch(ctx, mock.Anything, types.ExecStartCheck{}).Return(types.HijackedResponse{
		Reader: reader,
	}, nil)
	err := sandboxClusterExec(ctx, []string{}, cmdCtx)
	assert.NotNil(t, err)
}
