package create

import (
	cmdcore "github.com/flyteorg/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	createCmdShort = `Create various flyte resources including tasks/workflows/launchplans/executions/project.`
	createCmdLong  = `
The create command creates flyte resource, i.e project:
::

 flytectl create project --file project.yaml 
`
)

// RemoteCreateCommand will return create flyte resource commands
func RemoteCreateCommand() *cobra.Command {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: createCmdShort,
		Long:  createCmdLong,
	}
	createResourcesFuncs := map[string]cmdcore.CommandEntry{
		"project": {CmdFunc: createProjectsCommand, Aliases: []string{"projects"}, ProjectDomainNotRequired: true, PFlagProvider: projectConfig, Short: projectShort,
			Long: projectLong},
		"execution": {CmdFunc: createExecutionCommand, Aliases: []string{"executions"}, ProjectDomainNotRequired: false, PFlagProvider: executionConfig, Short: executionShort,
			Long: executionLong},
	}
	cmdcore.AddCommands(createCmd, createResourcesFuncs)
	return createCmd
}
