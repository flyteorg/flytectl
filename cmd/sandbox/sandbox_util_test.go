package sandbox

import (
	"io/ioutil"
	"os"
	"testing"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"

	"github.com/stretchr/testify/assert"
)

func TestConfigCleanup(t *testing.T) {
	_, err := os.Stat(f.FilePathJoin(f.UserHomeDir(), ".flyte"))
	if os.IsNotExist(err) {
		_ = os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte"), 0755)
	}
	_ = ioutil.WriteFile(FLYTECTLCONFIG, []byte("string"), 0600)
	_ = ioutil.WriteFile(KUBECONFIG, []byte("string"), 0600)

	err = ConfigCleanup()
	assert.Nil(t, err)

	_, err = os.Stat(FLYTECTLCONFIG)
	check := os.IsNotExist(err)
	assert.Equal(t, check, true)

	_, err = os.Stat(KUBECONFIG)
	check = os.IsNotExist(err)
	assert.Equal(t, check, true)
	_ = ConfigCleanup()
}

func TestSetupFlytectlConfig(t *testing.T) {
	_, err := os.Stat(f.FilePathJoin(f.UserHomeDir(), ".flyte"))
	if os.IsNotExist(err) {
		_ = os.MkdirAll(f.FilePathJoin(f.UserHomeDir(), ".flyte"), 0755)
	}
	err = SetupFlytectlConfig()
	assert.Nil(t, err)
	_, err = os.Stat(FLYTECTLCONFIG)
	assert.Nil(t, err)
	check := os.IsNotExist(err)
	assert.Equal(t, check, false)
	_ = ConfigCleanup()
}
