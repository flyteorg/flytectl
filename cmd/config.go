package cmd

import (
	"io"
	"os"

	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"

	f "github.com/flyteorg/flytectl/pkg/filesystemutils"
	"github.com/flyteorg/flytectl/pkg/util"
	"github.com/spf13/cobra"
)

var configFilePath = f.FilePathJoin(f.UserHomeDir(), ".flyte", "config.yaml")

// configCmd represents the config init command
var configCmd = &cobra.Command{
	Use:   "init-config",
	Short: "init-config flytectl config",
	Long: `
init-config will create flytectl config in flyte directory i.e ~/.flyte/config.yaml 
::

 bin/flytectl init-config
Usage`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initFlytectlConfig(os.Stdin)
	},
}

func initFlytectlConfig(reader io.Reader) error {
	if err := util.SetupFlyteDir(); err != nil {
		return err
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return util.WriteIntoFile([]byte(util.ConfigTemplate), f.FilePathJoin(f.UserHomeDir(), ".flyte", "config.yaml"))
	}

	if cmdUtil.AskForConfirmation("Are you sure ? It will overwrite the default config ~/.flyte/config.yaml", reader) {
		return util.WriteIntoFile([]byte(util.ConfigTemplate), f.FilePathJoin(f.UserHomeDir(), ".flyte", "config.yaml"))
	}
	return nil
}
