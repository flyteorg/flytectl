package init

import (
	"context"
	"fmt"
	"io"
	"os"

	initConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/init"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
	"github.com/flyteorg/flytectl/pkg/util"
)

const (
	initConfigCmdShort = "Teardown will cleanup the sandbox environment"
	initConfigCmdLong  = `

Usage
`
)

func configInitFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	return initFlytectlConfig(os.Stdin)
}

func initFlytectlConfig(reader io.Reader) error {
	if err := util.SetupFlyteDir(); err != nil {
		return err
	}
	spec := util.ConfigTemplateSpec{
		Host:     "dns:///localhost:30081",
		Insecure: true,
	}
	configTemplate := util.ConfigTemplate + util.StorageTemplate

	if len(initConfig.DefaultConfig.Host) > 0 {
		spec.Host = fmt.Sprintf("dns:///%v", initConfig.DefaultConfig.Host)
		spec.Insecure = false
		configTemplate = util.ConfigTemplate
	}

	if _, err := os.Stat(util.ConfigFile); os.IsNotExist(err) {
		return util.SetupConfig(configTemplate, util.ConfigFile, spec)
	}

	if cmdUtil.AskForConfirmation("Are you sure ? It will overwrite the default config ~/.flyte/config.yaml", reader) {
		if err := os.Remove(util.ConfigFile); err != nil {
			return err
		}
		return util.SetupConfig(configTemplate, util.ConfigFile, spec)
	}
	return nil
}
