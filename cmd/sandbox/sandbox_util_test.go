package sandbox

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	sandboxConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/sandbox"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	u "github.com/flyteorg/flytectl/cmd/testutils"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"

	"github.com/stretchr/testify/assert"
)

var (
	cmdCtx cmdCore.CommandContext
)

func cleanup(client *client.Client) error {
	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return err
	}
	for _, v := range containers {
		if strings.Contains(v.Names[0], flyteSandboxClusterName) {
			if err := client.ContainerRemove(context.Background(), v.ID, types.ContainerRemoveOptions{
				Force: true,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func setupSandbox() {
	mockAdminClient := u.MockClient
	cmdCtx = cmdCore.NewCommandContext(mockAdminClient, u.MockOutStream)
	_ = setupFlytectlConfig()
}

func TestConfigCleanup(t *testing.T) {
	_, err := os.Stat(f.FilePathJoin(f.UserHomeDir(), ".flyte"))
	if os.IsNotExist(err) {
		_ = os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte"), 0755)
	}
	_ = ioutil.WriteFile(FlytectlConfig, []byte("string"), 0600)
	_ = ioutil.WriteFile(Kubeconfig, []byte("string"), 0600)

	err = configCleanup()
	assert.Nil(t, err)

	_, err = os.Stat(FlytectlConfig)
	check := os.IsNotExist(err)
	assert.Equal(t, check, true)

	_, err = os.Stat(Kubeconfig)
	check = os.IsNotExist(err)
	assert.Equal(t, check, true)
	_ = configCleanup()
}

func TestSetupFlytectlConfig(t *testing.T) {
	_, err := os.Stat(f.FilePathJoin(f.UserHomeDir(), ".flyte"))
	if os.IsNotExist(err) {
		_ = os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte"), 0755)
	}
	err = setupFlytectlConfig()
	assert.Nil(t, err)
	_, err = os.Stat(FlytectlConfig)
	assert.Nil(t, err)
	check := os.IsNotExist(err)
	assert.Equal(t, check, false)
	_ = configCleanup()
}

func TestTearDownSandbox(t *testing.T) {
	setupSandbox()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	err := teardownSandboxCluster(context.Background(), []string{}, cmdCtx)
	assert.Nil(t, err)
	assert.Nil(t, cleanup(cli))

	_ = startSandboxCluster(context.Background(), []string{}, cmdCtx)
	err = teardownSandboxCluster(context.Background(), []string{}, cmdCtx)
	assert.Nil(t, err)
}

func TestStartSandboxErr(t *testing.T) {
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	assert.Nil(t, cleanup(cli))
	setupSandbox()
	volumes = []mount.Mount{}
	sandboxConfig.DefaultConfig.SnacksRepo = "/tmp"
	err := startSandboxCluster(context.Background(), []string{}, cmdCtx)
	assert.Nil(t, err)
}

func TestStartContainer(t *testing.T) {
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	assert.Nil(t, cleanup(cli))
	ID, err := startContainer(cli, []mount.Mount{})
	assert.Nil(t, err)
	err = readLogs(cli, ID, "Starting Docker daemon...")
	assert.Nil(t, err)
}
