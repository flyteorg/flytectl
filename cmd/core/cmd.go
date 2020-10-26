package cmdcore

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lyft/flytestdlib/logger"


	"github.com/lyft/flyteidl/clients/go/admin"
	"github.com/spf13/cobra"

	"github.com/lyft/flytectl/cmd/config"
)

type CustomFlags struct {
	P interface{}
	Name string
	Shorthand string
	Value interface{}
	Usage string
}

type CommandEntry struct {
	ProjectDomainNotRequired bool
	CmdFunc                  CommandFunc
	Aliases                  []string
	CustomFlags              []CustomFlags
	Subcommand               map[string]CommandEntry
}

func AddCommands(rootCmd *cobra.Command, cmdFuncs map[string]CommandEntry) {
	for resource, cmdEntry := range cmdFuncs {
		cmd := &cobra.Command{
			Use:     resource,
			Short:   fmt.Sprintf("Retrieves %v resources.", resource),
			Aliases: cmdEntry.Aliases,
			RunE:    generateCommandFunc(cmdEntry),
		}
		ctx := context.Background()
		for _, f := range cmdEntry.CustomFlags {
			switch f.P.(type) {
			case map[string]string:
				data, err := json.Marshal(f.P);
				if err != nil {
					logger.Debug(ctx, "Err : %v", err)
				}
				var value map[string]string
				err = json.Unmarshal(data, value);
				if err != nil {
					logger.Debug(ctx, "Err : %v", err)
				}
				valueData, err := json.Marshal(f.Value);
				if err != nil {
					logger.Debug(ctx, "Err : %v", err)
				}
				var valueUnmarshelData map[string]string
				err = json.Unmarshal(valueData, valueUnmarshelData);
				if err != nil {
					logger.Debug(ctx, "Err : %v", err)
				}
				rootCmd.PersistentFlags().StringToStringVar(&value, f.Name, valueUnmarshelData, f.Usage)
				break;
			case string:
				value := fmt.Sprintf("%v", f.P)
				rootCmd.PersistentFlags().StringVarP(&value, f.Name, f.Shorthand, fmt.Sprintf("%v", f.Value), f.Usage)
				break;
			default:
				fmt.Printf("I don't know about type %T!\n", f.P)

			}
			rootCmd.AddCommand(cmd)
			AddSubCommands(cmd,cmdEntry.Subcommand)
		}
	}
}

func AddSubCommands(rootCmd *cobra.Command, cmdFuncs map[string]CommandEntry) {
	for resource, cmdEntry := range cmdFuncs {
		cmd := &cobra.Command{
			Use:     resource,
			Short:   fmt.Sprintf("Retrieves %v sub resources.", resource),
			Aliases: cmdEntry.Aliases,
			RunE:    generateCommandFunc(cmdEntry),
		}
		ctx := context.Background()
		for _, f := range cmdEntry.CustomFlags {
			switch f.P.(type) {
			case map[string]string:
				data, err := json.Marshal(f.P);
				if err != nil {
					logger.Debug(ctx, "Err : %v", err)
				}
				var value map[string]string
				err = json.Unmarshal(data, value);
				if err != nil {
					logger.Debug(ctx, "Err : %v", err)
				}
				valueData, err := json.Marshal(f.Value);
				if err != nil {
					logger.Debug(ctx, "Err : %v", err)
				}
				var valueUnmarshelData map[string]string
				err = json.Unmarshal(valueData, valueUnmarshelData);
				if err != nil {
					logger.Debug(ctx, "Err : %v", err)
				}
				rootCmd.PersistentFlags().StringToStringVar(&value, f.Name, valueUnmarshelData, f.Usage)
				break;
			case string:
				value := fmt.Sprintf("%v", f.P)
				rootCmd.PersistentFlags().StringVarP(&value, f.Name, f.Shorthand, fmt.Sprintf("%v", f.Value), f.Usage)
				break;
			default:
				fmt.Printf("I don't know about type %T!\n", f.P)

			}
			rootCmd.AddCommand(cmd)
		}
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
