package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testVersion = "v0.1.20"

func TestWriteIntoFile(t *testing.T) {
	t.Run("Successfully write into a file", func(t *testing.T) {
		err := WriteIntoFile([]byte("data"), "version.yaml")
		assert.Nil(t, err)
	})
	t.Run("Error in writing file", func(t *testing.T) {
		err := WriteIntoFile([]byte("data"), "/githubtest/version.yaml")
		assert.NotNil(t, err)
	})
}

func TestSetupFlyteDir(t *testing.T) {
	assert.Nil(t, SetupFlyteDir())
}

func TestIsVersionGreaterThan(t *testing.T) {
	t.Run("Compare flytectl version when upgrade available", func(t *testing.T) {
		_, err := IsVersionGreaterThan("v1.1.21", testVersion)
		assert.Nil(t, err)
	})
	t.Run("Compare flytectl version greater then", func(t *testing.T) {
		ok, err := IsVersionGreaterThan("v1.1.21", testVersion)
		assert.Nil(t, err)
		assert.Equal(t, true, ok)
	})
	t.Run("Compare flytectl version smaller then", func(t *testing.T) {
		ok, err := IsVersionGreaterThan("v0.1.19", testVersion)
		assert.Nil(t, err)
		assert.Equal(t, false, ok)
	})
	t.Run("Compare flytectl version", func(t *testing.T) {
		_, err := IsVersionGreaterThan(testVersion, testVersion)
		assert.Nil(t, err)
	})
	t.Run("Error in compare flytectl version", func(t *testing.T) {
		_, err := IsVersionGreaterThan("vvvvvvvv", testVersion)
		assert.NotNil(t, err)
	})
	t.Run("Error in compare flytectl version", func(t *testing.T) {
		_, err := IsVersionGreaterThan(testVersion, "vvvvvvvv")
		assert.NotNil(t, err)
	})
}

func TestGetLatestRelease(t *testing.T) {
	release, err := GetLatestVersion("flyte")
	assert.Nil(t, err)
	assert.Equal(t, true, strings.HasPrefix(release.GetTagName(), "v"))
}

func TestCheckVersionExist(t *testing.T) {
	t.Run("Invalid Tag", func(t *testing.T) {
		_, err := CheckVersionExist("v100.0.0", "flyte")
		assert.NotNil(t, err)
	})
	t.Run("Valid Tag", func(t *testing.T) {
		release, err := CheckVersionExist("v0.15.0", "flyte")
		assert.Nil(t, err)
		assert.Equal(t, true, strings.HasPrefix(release.GetTagName(), "v"))
	})
}
