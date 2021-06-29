package cmd

import (
	"io"
	"os"

	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
	"github.com/flyteorg/flytectl/pkg/docker"
	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	"github.com/flyteorg/flytectl/pkg/util"
	"github.com/spf13/cobra"
)

// configCmd represents the config init command
var configCmd = &cobra.Command{
	Use:   "init-config",
	Short: "Init flytectl config",
	Long:  `config will create a flytectl config in defined path`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initFlytectlConfig(os.Stdin)
	},
}

func initFlytectlConfig(reader io.Reader) error {
	if err := docker.SetupFlyteDir(); err != nil {
		return err
	}
	if cmdUtil.AskForConfirmation("Are you sure ? It will overwrite the default config from ~/.flyte/config.yaml", reader) {
		return util.WriteIntoFile([]byte(util.ConfigTemplate), f.FilePathJoin(f.UserHomeDir(), ".flyte", "config.yaml"))
	}
	return nil
}
