package sandbox

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	u "github.com/flyteorg/flytectl/cmd/testutils"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"

	"github.com/stretchr/testify/assert"
)

var (
	cmdCtx cmdCore.CommandContext
)

func setupSandbox() {
	mockAdminClient := u.MockClient
	cmdCtx = cmdCore.NewCommandContext(mockAdminClient, u.MockOutStream)
}
func TestConfigCleanup(t *testing.T) {
	_, err := os.Stat(f.FilePathJoin(f.UserHomeDir(), ".flyte"))
	if os.IsNotExist(err) {
		_ = os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte"), 0755)
	}
	_ = ioutil.WriteFile(FLYTECTLCONFIG, []byte("string"), 0600)
	_ = ioutil.WriteFile(KUBECONFIG, []byte("string"), 0600)

	err = configCleanup()
	assert.Nil(t, err)

	_, err = os.Stat(FLYTECTLCONFIG)
	check := os.IsNotExist(err)
	assert.Equal(t, check, true)

	_, err = os.Stat(KUBECONFIG)
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
	_, err = os.Stat(FLYTECTLCONFIG)
	assert.Nil(t, err)
	check := os.IsNotExist(err)
	assert.Equal(t, check, false)
	_ = configCleanup()
}

func cleanup(client *client.Client) error {
	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return err
	}
	for _, v := range containers {
		if strings.Contains(v.Names[0], SandboxClusterName) {
			if err := client.ContainerRemove(context.Background(), v.ID, types.ContainerRemoveOptions{}); err != nil {
				return err
			}
		}
	}
	return nil
}

func TestGetSandbox(t *testing.T) {
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	assert.Nil(t, cleanup(cli))
	container := getSandbox(cli)
	assert.Nil(t, container)
	assert.Nil(t, cleanup(cli))

}

func TestTearDownSandbox(t *testing.T) {
	setupSandbox()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	err := teardownSandboxCluster(context.Background(), []string{}, cmdCtx)
	assert.Nil(t, err)
	assert.Nil(t, cleanup(cli))
}

func TestStartContainer(t *testing.T) {
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	assert.Nil(t, cleanup(cli))
	id, err := startContainer(cli)
	assert.Nil(t, err)
	assert.Greater(t, len(id), 0)
}
