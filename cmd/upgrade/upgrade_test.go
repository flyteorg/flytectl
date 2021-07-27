package upgrade

import (
	"io"
	"sort"
	"strings"
	"testing"

	"github.com/flyteorg/flytectl/pkg/filesystemutils"

	"github.com/flyteorg/flytectl/pkg/util"

	"github.com/flyteorg/flyteidl/clients/go/admin/mocks"
	stdlibversion "github.com/flyteorg/flytestdlib/version"

	"context"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestUpgradeCommand(t *testing.T) {
	rootCmd := &cobra.Command{
		Long:              "flytectl is CLI tool written in go to interact with flyteadmin service",
		Short:             "flyetcl CLI tool",
		Use:               "flytectl",
		DisableAutoGenTag: true,
	}
	upgradeCmd := SelfUpgrade(rootCmd)
	cmdCore.AddCommands(rootCmd, upgradeCmd)
	assert.Equal(t, len(rootCmd.Commands()), 1)
	cmdNouns := rootCmd.Commands()
	// Sort by Use value.
	sort.Slice(cmdNouns, func(i, j int) bool {
		return cmdNouns[i].Use < cmdNouns[j].Use
	})

	assert.Equal(t, cmdNouns[0].Use, "upgrade")
	assert.Equal(t, cmdNouns[0].Short, upgradeCmdShort)
	assert.Equal(t, cmdNouns[0].Long, upgradeCmdLong)
}

func TestUpgrade(t *testing.T) {
	t.Run("Successful upgrade", func(t *testing.T) {
		release, err := util.CheckVersionExist("v0.2.14", "flytectl")
		if err != nil {
			t.Error(err)
		}
		assert.Nil(t, upgrade(strings.NewReader("y"), release, "/tmp/flytectl", "linux"))
		assert.Nil(t, upgrade(strings.NewReader("n"), release, "/tmp/flytectl", "linux"))
		assert.Nil(t, upgrade(strings.NewReader("n"), release, "/tmp/flytectl", "windows"))
	})
}

func TestSelfUpgrade(t *testing.T) {
	ext = filesystemutils.FilePathJoin("/tmp/test")
	_ = util.WriteIntoFile([]byte(""), ext)

	t.Run("Successful upgrade", func(t *testing.T) {
		ctx := context.Background()
		var args []string
		mockClient := new(mocks.AdminServiceClient)
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
		stdlibversion.Build = ""
		stdlibversion.BuildTime = ""
		stdlibversion.Version = "v0.2.10"

		assert.Nil(t, selfUpgrade(ctx, args, cmdCtx))
	})
}

func TestSelfUpgradeRollback(t *testing.T) {
	ext = filesystemutils.FilePathJoin("/tmp/test")
	_ = util.WriteIntoFile([]byte(""), ext)

	t.Run("Successful upgrade", func(t *testing.T) {
		ctx := context.Background()
		var args = []string{"rollback"}
		mockClient := new(mocks.AdminServiceClient)
		mockOutStream := new(io.Writer)
		cmdCtx := cmdCore.NewCommandContext(mockClient, *mockOutStream)
		stdlibversion.Build = ""
		stdlibversion.BuildTime = ""
		stdlibversion.Version = "v0.2.10"

		assert.Nil(t, selfUpgrade(ctx, args, cmdCtx))
	})
}
