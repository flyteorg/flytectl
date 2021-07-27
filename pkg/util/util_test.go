package util

import (
	"fmt"
	"strings"
	"testing"
	"time"

	stdlibversion "github.com/flyteorg/flytestdlib/version"

	"github.com/stretchr/testify/assert"
)

const testVersion = "v0.1.20"

func TestWriteIntoFile(t *testing.T) {
	t.Run("Successfully write into a file", func(t *testing.T) {
		err := WriteIntoFile([]byte(""), "version.yaml")
		assert.Nil(t, err)
	})
	t.Run("Error in writing file", func(t *testing.T) {
		err := WriteIntoFile([]byte(""), "version.yaml")
		assert.Nil(t, err)
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

func TestProgressBarForFlyteStatus(t *testing.T) {
	t.Run("Progress bar success", func(t *testing.T) {
		count := make(chan int64)
		go func() {
			time.Sleep(1 * time.Second)
			count <- 1
		}()
		ProgressBarForFlyteStatus(1, count, "")
	})
}

func TestGetLatestVersion(t *testing.T) {
	t.Run("Get latest release with wrong url", func(t *testing.T) {
		_, err := GetLatestVersion("fl")
		assert.NotNil(t, err)
	})
	t.Run("Get latest release", func(t *testing.T) {
		_, err := GetLatestVersion("flytectl")
		assert.Nil(t, err)
	})
}

func TestDetectNewVersion(t *testing.T) {
	stdlibversion.Version = "v0.2.10"
	message, err := GetUpgradeMessage("darwin")
	fmt.Println(message)
	assert.Nil(t, err)
	assert.Equal(t, 177, len(message))
	stdlibversion.Version = "v0.2.0"
	message, err = GetUpgradeMessage("darwin")
	assert.Nil(t, err)
	assert.Equal(t, 176, len(message))
	stdlibversion.Version = "v100.0.0"
	message, err = GetUpgradeMessage("darwin")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(message))
	stdlibversion.Version = "v0"
	message, err = GetUpgradeMessage("darwin")
	assert.Nil(t, err)
	assert.Equal(t, 172, len(message))
	message, err = GetUpgradeMessage("linux")
	assert.Nil(t, err)
	assert.Equal(t, 152, len(message))
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
