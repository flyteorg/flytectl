package init

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/manifoldco/promptui"

	initConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/init"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
	"github.com/flyteorg/flytectl/pkg/util"
)

var prompt = promptui.Select{
	Label: "Select Storage Provider",
	Items: []string{"S3", "GCS"},
}

func configInitFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	return initFlytectlConfig(os.Stdin)
}

func initFlytectlConfig(reader io.Reader) error {

	templateValues := util.ConfigTemplateValuesSpec{
		Host:     "dns:///localhost:30081",
		Insecure: initConfig.DefaultConfig.Insecure,
		Template: util.GetSandboxTemplate(),
	}

	if len(initConfig.DefaultConfig.Host) > 0 {
		templateValues.Host = fmt.Sprintf("dns:///%v", initConfig.DefaultConfig.Host)
		templateValues.Template = util.GetAWSCloudTemplate()
		_, result, err := prompt.Run()
		if err != nil {
			return err
		}
		if result == "GCS" {
			templateValues.Template = util.GetGoogleCloudTemplate()
		}
	}
	var _err error
	if _, err := os.Stat(util.ConfigFile); os.IsNotExist(err) {
		_err = util.SetupConfig(util.ConfigFile, templateValues)
	} else {
		if cmdUtil.AskForConfirmation(fmt.Sprintf("This action will overwrite an existing config file at [%s]. Do you want to continue?", util.ConfigFile), reader) {
			if err := os.Remove(util.ConfigFile); err != nil {
				return err
			}
			_err = util.SetupConfig(util.ConfigFile, templateValues)
		}
	}

	if len(initConfig.DefaultConfig.Host) > 0 {
		fmt.Println("Init flytectl config for remote cluster, Please update your storage config in ~/.flyte/config.yaml. Learn more about the config here https://docs.flyte.org/projects/flytectl/en/latest/index.html#configure")
	}
	return _err
}
