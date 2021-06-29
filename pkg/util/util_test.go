package util

import (
	"io/ioutil"
	"os"
	"testing"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"

	"github.com/stretchr/testify/assert"
)

const flytectlReleaseURL = "/repos/flyteorg/flytectl/releases/latest"
const baseURL = "https://api.github.com"
const wrongBaseURL = "htts://api.github.com"

func TestGetRequest(t *testing.T) {
	t.Run("Get request with 200", func(t *testing.T) {
		_, err := GetRequest(baseURL, flytectlReleaseURL)
		assert.Nil(t, err)
	})
	t.Run("Get request with 200", func(t *testing.T) {
		_, err := GetRequest(wrongBaseURL, flytectlReleaseURL)
		assert.NotNil(t, err)
	})
}

func TestParseGithubTag(t *testing.T) {
	t.Run("Parse Github tag with success", func(t *testing.T) {
		data, err := GetRequest(baseURL, flytectlReleaseURL)
		assert.Nil(t, err)
		tag, err := ParseGithubTag(data)
		assert.Nil(t, err)
		assert.Contains(t, tag, "v")
	})
	t.Run("Get request with 200", func(t *testing.T) {
		_, err := ParseGithubTag([]byte("string"))
		assert.NotNil(t, err)
	})
}

func TestWriteIntoFile(t *testing.T) {
	t.Run("Successfully write into a file", func(t *testing.T) {
		data, err := GetRequest(baseURL, flytectlReleaseURL)
		assert.Nil(t, err)
		err = WriteIntoFile(data, "version.yaml")
		assert.Nil(t, err)
	})
	t.Run("Error in writing file", func(t *testing.T) {
		data, err := GetRequest(baseURL, flytectlReleaseURL)
		assert.Nil(t, err)
		err = WriteIntoFile(data, "/githubtest/version.yaml")
		assert.NotNil(t, err)
	})
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
	err = SetupConfig()
	assert.Nil(t, err)
	_, err = os.Stat(FlytectlConfig)
	assert.Nil(t, err)
	check := os.IsNotExist(err)
	assert.Equal(t, check, false)
	_ = ConfigCleanup()
}
