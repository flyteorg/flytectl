package cmd

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/cmd/create"
	"github.com/flyteorg/flytectl/cmd/delete"
	"github.com/flyteorg/flytectl/cmd/get"
	"github.com/flyteorg/flytectl/cmd/register"
	"github.com/flyteorg/flytectl/cmd/update"
	"github.com/flyteorg/flytectl/cmd/version"
	"github.com/flyteorg/flytectl/pkg/printer"
	stdConfig "github.com/flyteorg/flytestdlib/config"
	"github.com/flyteorg/flytestdlib/config/viper"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	cfgFile        string
	configAccessor = viper.NewAccessor(stdConfig.Options{StrictMode: true})
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		PersistentPreRunE: initConfig,
		Long:              "flytectl is CLI tool written in go to interact with flyteadmin service",
		Short:             "flyetcl CLI tool",
		Use:               "flytectl",
		DisableAutoGenTag: true,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/config.yaml)")

	configAccessor.InitializePflags(rootCmd.PersistentFlags())

	// Due to https://github.com/lyft/flyte/issues/341, project flag will have to be specified as
	// --root.project, this adds a convenience on top to allow --project to be used
	rootCmd.PersistentFlags().StringVarP(&(config.GetConfig().Project), "project", "p", "", "Specifies the Flyte project.")
	rootCmd.PersistentFlags().StringVarP(&(config.GetConfig().Domain), "domain", "d", "", "Specifies the Flyte project's domain.")
	rootCmd.PersistentFlags().StringVarP(&(config.GetConfig().Output), "output", "o", printer.OutputFormatTABLE.String(), fmt.Sprintf("Specifies the output type - supported formats %s", printer.OutputFormats()))
	rootCmd.PersistentFlags().StringVarP(&(config.GetConfig().FieldSelector), "field-selector", "f", "", "Specifies the Field selector")
	rootCmd.PersistentFlags().StringVarP(&(config.GetConfig().SortBy), "sort-by", "s", "", "Specifies which field to sort results ")
	rootCmd.PersistentFlags().Int32Var(&(config.GetConfig().Limit), "limit", config.DefaultLimit, "Specifies the limit of results")
	rootCmd.PersistentFlags().BoolVar(&(config.GetConfig().Asc), "asc", false, "Specifies the sorting order. By default flytectl sort result in descending order")
	rootCmd.AddCommand(viper.GetConfigCommand())
	rootCmd.AddCommand(get.CreateGetCommand())
	rootCmd.AddCommand(create.RemoteCreateCommand())
	rootCmd.AddCommand(update.CreateUpdateCommand())
	rootCmd.AddCommand(register.RemoteRegisterCommand())
	rootCmd.AddCommand(delete.RemoteDeleteCommand())
	// Added version command
	versioncmd := version.GetVersionCommand(rootCmd)
	cmdCore.AddCommands(rootCmd, versioncmd)

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

func GenerateDocs() error {
	rootCmd := newRootCmd()
	err := GenReSTTree(rootCmd, "gen")
	if err != nil {
		logrus.Fatal(err)
		return err
	}
	return nil
}

func GenReSTTree(cmd *cobra.Command, dir string) error {
	emptyStr := func(s string) string { return "" }
	// Sphinx cross-referencing format
	linkHandler := func(name, ref string) string {
		return fmt.Sprintf(":doc:`%s`", ref)
	}
	return doc.GenReSTTreeCustom(cmd, dir, emptyStr, linkHandler)
}

func ExecuteCmd() error {
	return newRootCmd().Execute()
}
