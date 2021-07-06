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
	initConfigCmdShort = "Used for generating config template."
	initConfigCmdLong  = `init config will create a config in flyte directory i.e ~/.flyte
Generate sandbox config.
	
::

 bin/flytectl init config 

Generate remote cluster config. 
	
::

 bin/flytectl init config --host="flyte.myexample.com"
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
	var _err error
	if _, err := os.Stat(util.ConfigFile); os.IsNotExist(err) {
		_err = util.SetupConfig(configTemplate, util.ConfigFile, spec)
	} else {
		if cmdUtil.AskForConfirmation("Are you sure ? It will overwrite the default config ~/.flyte/config.yaml", reader) {
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
