package configuration

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/flyteorg/flytestdlib/logger"

	"github.com/flyteorg/flytectl/pkg/util"

	"github.com/flyteorg/flytectl/pkg/configutil"

	initConfig "github.com/flyteorg/flytectl/cmd/config/subcommand/config"
	cmdcore "github.com/flyteorg/flytectl/cmd/core"
	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
	"github.com/flyteorg/flytestdlib/config/viper"
	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when docs are generated using sphinx.
const (
	initCmdShort = `Generates flytectl config file in the user's home directory.`
	initCmdLong  = `Creates a flytectl config file in flyte directory i.e ~/.flyte
	
It generates a flytectl config file in the flyte directory, i.e ~/.flyte.
It is sandbox config known as 'Flyte Sandbox'. Flyte Sandbox is a complete standalone minimal environment to run Flyte. 
Read more about sandbox `here. <https://docs.flyte.org/en/latest/deployment/sandbox.html>`_
::

 bin/flytectl configuration config 

A remote cluster config can be generated with the below command: 
Read more about the remote deployment 'here. <https://docs.flyte.org/en/latest/deployment/index.html>'_	
::

 bin/flytectl configuration config --host=flyte.myexample.com
	
A flytectl config can be generated using a storage provider. It can be generated using the below command:
::

 bin/flytectl configuration config --host=flyte.myexample.com --storage
`
)

var prompt = promptui.Select{
	Label: "Select Storage Provider",
	Items: []string{"S3", "GCS"},
}

// CreateConfigCommand will return configuration command
func CreateConfigCommand() *cobra.Command {
	configCmd := viper.GetConfigCommand()

	getResourcesFuncs := map[string]cmdcore.CommandEntry{
		"init": {CmdFunc: configInitFunc, Aliases: []string{""}, ProjectDomainNotRequired: true,
			Short: initCmdShort,
			Long:  initCmdLong, PFlagProvider: initConfig.DefaultConfig},
	}

	cmdcore.AddCommands(configCmd, getResourcesFuncs)
	return configCmd
}

func configInitFunc(ctx context.Context, args []string, cmdCtx cmdcore.CommandContext) error {
	return initFlytectlConfig(ctx, os.Stdin)
}

func initFlytectlConfig(ctx context.Context, reader io.Reader) error {

	if err := util.SetupFlyteDir(); err != nil {
		return err
	}

	templateValues := configutil.ConfigTemplateSpec{
		Host:     "dns:///localhost:30081",
		Insecure: initConfig.DefaultConfig.Insecure,
	}
	templateStr := configutil.GetSandboxTemplate()

	if len(initConfig.DefaultConfig.Host) > 0 {
		templateValues.Host = fmt.Sprintf("dns:///%v", initConfig.DefaultConfig.Host)
		templateStr = configutil.AdminConfigTemplate
		if initConfig.DefaultConfig.Storage {
			templateStr = configutil.GetAWSCloudTemplate()
			_, result, err := prompt.Run()
			if err != nil {
				return err
			}
			if strings.ToUpper(result) == "GCS" {
				templateStr = configutil.GetGoogleCloudTemplate()
			}
		} else {
			logger.Infof(ctx, "Init flytectl config for remote cluster, Please update your storage config in %s. Learn more about the config here https://docs.flyte.org/projects/flytectl/en/latest/index.html#configure", configutil.ConfigFile)
		}
	}
	var _err error
	if _, err := os.Stat(configutil.ConfigFile); os.IsNotExist(err) {
		_err = configutil.SetupConfig(configutil.ConfigFile, templateStr, templateValues)
	} else {
		if cmdUtil.AskForConfirmation(fmt.Sprintf("This action will overwrite an existing config file at [%s]. Do you want to continue?", configutil.ConfigFile), reader) {
			if err := os.Remove(configutil.ConfigFile); err != nil {
				return err
			}
			_err = configutil.SetupConfig(configutil.ConfigFile, templateStr, templateValues)
		}
	}
	if _err != nil {
		return _err
	}
	fmt.Printf("Init flytectl config file at [%s]", configutil.ConfigFile)
	return nil
}
