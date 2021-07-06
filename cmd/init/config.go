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

func configInitFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	return initFlytectlConfig(os.Stdin)
}

func initFlytectlConfig(reader io.Reader) error {
	if err := util.SetupFlyteDir(); err != nil {
		return err
	}
	spec := util.ConfigTemplateValuesSpec{
		Host:     "dns:///localhost:30081",
		Insecure: true,
	}
	configTemplate := util.GetSandboxTemplate()

	if len(initConfig.DefaultConfig.Host) > 0 {
		spec.Host = fmt.Sprintf("dns:///%v", initConfig.DefaultConfig.Host)
		spec.Insecure = true
		configTemplate = util.AdminConfigTemplate
	}
	var _err error
	if _, err := os.Stat(util.ConfigFile); os.IsNotExist(err) {
		_err = util.SetupConfig(configTemplate, util.ConfigFile, spec)
	} else {
		if cmdUtil.AskForConfirmation(fmt.Sprintf("Are you sure ? It will overwrite the default config %v", util.ConfigFile), reader) {
			if err := os.Remove(util.ConfigFile); err != nil {
				return err
			}
			_err = util.SetupConfig(configTemplate, util.ConfigFile, spec)
		}
	}

	if len(initConfig.DefaultConfig.Host) > 0 {
		fmt.Println("Init flytectl config for remote cluster, Please update your storage config in ~/.flyte/config.yaml. Learn more about the config here https://docs.flyte.org/projects/flytectl/en/latest/index.html#configure")
	}
	return _err
}
