package cmdcore

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"

	"github.com/lyft/flyteidl/clients/go/admin"
	"github.com/spf13/cobra"

	"github.com/lyft/flytectl/cmd/config"
)

type PFlagProvider interface {
	GetPFlagSet(prefix string) *pflag.FlagSet
}

type CommandEntry struct {
	ProjectDomainNotRequired bool
	CmdFunc                  CommandFunc
	Aliases                  []string
	PFlagProvider            PFlagProvider
}

func AddCommands(rootCmd *cobra.Command, cmdFuncs map[string]CommandEntry) {
	for resource, cmdEntry := range cmdFuncs {
		cmd := &cobra.Command{
			Use:     resource,
			Short:   fmt.Sprintf("Retrieves %v resources.", resource),
			Aliases: cmdEntry.Aliases,
			RunE:    generateCommandFunc(cmdEntry),
		}
		if cmdEntry.PFlagProvider != nil {
			cmd.Flags().AddFlagSet(cmdEntry.PFlagProvider.GetPFlagSet(""))
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
		return cmdEntry.CmdFunc(ctx, args, NewCommandContext(adminClient, cmd.OutOrStdout()))
	}
}
