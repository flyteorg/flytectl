package sandbox

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
	"github.com/flyteorg/flytectl/pkg/docker"
	"github.com/flyteorg/flytectl/pkg/docker/mocks"
	"github.com/moby/moby/pkg/stdcopy"
	"github.com/stretchr/testify/mock"
)

func getSrcBuffer(stdOutBytes, stdErrBytes []byte) (buffer *bytes.Buffer, err error) {
	buffer = new(bytes.Buffer)
	dstOut := stdcopy.NewStdWriter(buffer, stdcopy.Stdout)
	_, err = dstOut.Write(stdOutBytes)
	if err != nil {
		return
	}
	dstErr := stdcopy.NewStdWriter(buffer, stdcopy.Stderr)
	_, err = dstErr.Write(stdErrBytes)
	return
}

func TestSandboxClusterExec(t *testing.T) {
	mockDocker := &mocks.Docker{}
	mockOutStream := new(io.Writer)
	ctx := context.Background()
	cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
	c := docker.ExecConfig
	c.Cmd = []string{"ls"}
	stdOutBytes := []byte(strings.Repeat("o", docker.StartingBufLen))
	stdErrBytes := []byte(strings.Repeat("e", docker.StartingBufLen))
	buffer, err := getSrcBuffer(stdOutBytes, stdErrBytes)
	reader := bufio.NewReader(buffer)
	if err != nil {
		t.Fatal(err)
	}
	mockDocker.OnContainerExecCreateMatch(ctx, mock.Anything, c).Return(types.IDResponse{}, nil)
	mockDocker.OnContainerExecInspectMatch(ctx, mock.Anything).Return(types.ContainerExecInspect{}, nil)
	mockDocker.OnContainerExecAttachMatch(ctx, mock.Anything, types.ExecStartCheck{}).Return(types.HijackedResponse{
		Reader: reader,
	}, nil)
	err = sandboxClusterExec(ctx, []string{"ls -al"}, cmdCtx)
	assert.NotNil(t, err)
}

func TestSandboxClusterExecWithNoCommand(t *testing.T) {
	mockDocker := &mocks.Docker{}
	mockOutStream := new(io.Writer)
	ctx := context.Background()
	cmdCtx := cmdCore.NewCommandContext(nil, *mockOutStream)
	c := docker.ExecConfig
	c.Cmd = []string{"ls"}
	stdOutBytes := []byte(strings.Repeat("o", docker.StartingBufLen))
	stdErrBytes := []byte(strings.Repeat("e", docker.StartingBufLen))
	buffer, err := getSrcBuffer(stdOutBytes, stdErrBytes)
	reader := bufio.NewReader(buffer)
	if err != nil {
		t.Fatal(err)
	}
	mockDocker.OnContainerExecCreateMatch(ctx, mock.Anything, c).Return(types.IDResponse{}, nil)
	mockDocker.OnContainerExecInspectMatch(ctx, mock.Anything).Return(types.ContainerExecInspect{}, nil)
	mockDocker.OnContainerExecAttachMatch(ctx, mock.Anything, types.ExecStartCheck{}).Return(types.HijackedResponse{
		Reader: reader,
	}, nil)
	err = sandboxClusterExec(ctx, []string{}, cmdCtx)
	assert.NotNil(t, err)
}
