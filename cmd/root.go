package cmd

import (
	"context"

	"github.com/lyft/flytestdlib/config"
	"github.com/lyft/flytestdlib/config/viper"
	"github.com/spf13/cobra"
)

var (
	cfgFile        string
	configAccessor = viper.NewAccessor(config.Options{StrictMode: true})
)

type persistentFlags struct {
	Project *string
	Domain  *string
}

func newRootCmd() *cobra.Command {
	persistentFlags := persistentFlags{}
	rootCmd := &cobra.Command{
		PersistentPreRunE: initConfig,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/config.yaml)")

	configAccessor.InitializePflags(rootCmd.PersistentFlags())

	persistentFlags.Project = rootCmd.PersistentFlags().String("project", "", "Specifies the Flyte project.")
	persistentFlags.Domain = rootCmd.PersistentFlags().String("domain", "", "Specifies the Flyte project's domain.")

	rootCmd.AddCommand(newTimelineCmd(persistentFlags))
	return rootCmd
}

func initConfig(_ *cobra.Command, _ []string) error {
	configAccessor = viper.NewAccessor(config.Options{
		StrictMode:  true,
		SearchPaths: []string{cfgFile},
	})

	err := configAccessor.UpdateConfig(context.TODO())
	if err != nil {
		return err
	}

	return nil
}

func ExecuteCmd() error {
	return newRootCmd().Execute()
}
