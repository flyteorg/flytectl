package cmdcore

import (
	"context"
	"fmt"

	"github.com/lyft/flyteidl/clients/go/admin"
	"github.com/spf13/cobra"

	"github.com/lyft/flytectl/cmd/config"
)

type CustomFlags struct {
	P *string
	Name string
	Shorthand string
	Value string
	Usage string
}

type CommandEntry struct {
	ProjectDomainNotRequired bool
	CmdFunc                  CommandFunc
	Aliases                  []string
	CustomFlags              []CustomFlags
}

func AddCommands(rootCmd *cobra.Command, cmdFuncs map[string]CommandEntry) {
	for resource, cmdEntry := range cmdFuncs {
		cmd := &cobra.Command{
			Use:     resource,
			Short:   fmt.Sprintf("Retrieves %v resources.", resource),
			Aliases: cmdEntry.Aliases,
			RunE:    generateCommandFunc(cmdEntry),
		}
		for _, f := range cmdEntry.CustomFlags {
			rootCmd.PersistentFlags().StringVarP(f.P, f.Name, f.Shorthand,f.Value, f.Usage)
		}
		rootCmd.AddCommand(cmd)
	}
}

func generateCommandFunc(cmdEntry CommandEntry) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		if !cmdEntry.ProjectDomainNotRequired {
			if config.GetConfig().Project == "" {
				return fmt.Errorf("project and domain are required parameters")
			}
			if config.GetConfig().Domain == "" {
				return fmt.Errorf("project and domain are required parameters")
			}
		}
		if _, err := config.GetConfig().OutputFormat(); err != nil {
			return err
		}

		adminClient, err := admin.InitializeAdminClientFromConfig(ctx)
		if err != nil {
			return err
		}
		return cmdEntry.CmdFunc(ctx, args, CommandContext{
			out:         cmd.OutOrStdout(),
			adminClient: adminClient,
		})
	}
}
