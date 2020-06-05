package cmd

import (
	"context"
	"fmt"

	"github.com/lyft/flytectl/cmd/get"

	"github.com/lyft/flytectl/cmd/config"
	"github.com/lyft/flyteidl/clients/go/admin"

	stdConfig "github.com/lyft/flytestdlib/config"
	"github.com/lyft/flytestdlib/config/viper"
	"github.com/spf13/cobra"
)

var (
	cfgFile        string
	configAccessor = viper.NewAccessor(stdConfig.Options{StrictMode: true})
)

type CmdFunc func(ctx context.Context, args []string, cmdCtx CommandContext) error

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		PersistentPreRunE: initConfig,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/config.yaml)")

	configAccessor.InitializePflags(rootCmd.PersistentFlags())

	// Due to https://github.com/lyft/flyte/issues/341, project flag will have to be specified as
	// --root.project, this adds a convenience on top to allow --project to be used
	rootCmd.PersistentFlags().StringVarP(&(config.GetConfig().Project), "project", "p", "", "Specifies the Flyte project.")
	rootCmd.PersistentFlags().StringVarP(&(config.GetConfig().Domain), "domain", "d", "", "Specifies the Flyte project's domain.")

	rootCmd.AddCommand(viper.GetConfigCommand())
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(get.CreateGetCommand())
	config.GetConfig()

	return rootCmd
}

func initConfig(_ *cobra.Command, _ []string) error {
	configAccessor = viper.NewAccessor(stdConfig.Options{
		StrictMode:  true,
		SearchPaths: []string{cfgFile},
	})

	err := configAccessor.UpdateConfig(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

func AddCommands(rootCmd *cobra.Command, cmdFuncs map[string]CmdFunc) {
	for resource, getFunc := range cmdFuncs {
		cmd := &cobra.Command{
			Use:   resource,
			Short: fmt.Sprintf("Retrieves %v resources.", resource),
			RunE:  generateCommandFunc(getFunc),
		}

		rootCmd.AddCommand(cmd)
	}
}

func ExecuteCmd() error {
	return newRootCmd().Execute()
}

func generateCommandFunc(cmdFunc CmdFunc) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		adminClient, err := admin.InitializeAdminClientFromConfig(ctx)
		if err != nil {
			return err
		}

		return cmdFunc(ctx, args, CommandContext{
			out:         cmd.OutOrStdout(),
			adminClient: adminClient,
		})
	}
}
