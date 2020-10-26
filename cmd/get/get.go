package get

import (
	cmdcore "github.com/lyft/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// CreateGetCommand will return get command
func CreateGetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve various resource.",
	}

	getResourcesFuncs := map[string]cmdcore.CommandEntry{
		"project":    {CmdFunc: getProjectsFunc, Aliases: []string{"projects"}, ProjectDomainNotRequired: true},
		"task":       {CmdFunc: getTaskFunc, Aliases: []string{"tasks"}},
		"workflow":   {CmdFunc: getWorkflowFunc, Aliases: []string{"workflows"}},
		"launchplan": {CmdFunc: getLaunchPlanFunc, Aliases: []string{"launchplans"}},
		//"execution": {CmdFunc: getExecutionFunc, Aliases: []string{"executions"},CustomFlags : []cmdcore.CustomFlags{},Subcommand: map[string]cmdcore.CommandEntry{
		//	"node":       {CmdFunc: getExecutionFunc, Aliases: []string{""}},
		//	"task":       {CmdFunc: getExecutionFunc, Aliases: []string{""}},
		//}},
	}

	cmdcore.AddCommands(getCmd, getResourcesFuncs)

	return getCmd
}
