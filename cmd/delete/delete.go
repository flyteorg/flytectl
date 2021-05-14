package delete

import (
	"github.com/flyteorg/flytectl/cmd/config/subcommand/clusterresourceattribute"
	"github.com/flyteorg/flytectl/cmd/config/subcommand/executionclusterlabel"
	"github.com/flyteorg/flytectl/cmd/config/subcommand/executionqueueattribute"
	"github.com/flyteorg/flytectl/cmd/config/subcommand/taskresourceattribute"
	cmdcore "github.com/flyteorg/flytectl/cmd/core"

	"github.com/spf13/cobra"
)

// Long descriptions are whitespace sensitive when generating docs using sphinx.
const (
	deleteCmdShort = `Used for terminating/deleting various flyte resources including tasks/workflows/launchplans/executions/project.`
	deleteCmdLong  = `
Example Delete executions.
::

 bin/flytectl delete execution kxd1i72850  -d development  -p flytesnacks
`
)

// RemoteDeleteCommand will return delete command
func RemoteDeleteCommand() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: deleteCmdShort,
		Long:  deleteCmdLong,
	}
	terminateResourcesFuncs := map[string]cmdcore.CommandEntry{
		"execution": {CmdFunc: terminateExecutionFunc, Aliases: []string{"executions"}, Short: execCmdShort,
			Long: execCmdLong},
		"task-resource-attribute": {CmdFunc: deleteTaskResourceAttributes, Aliases: []string{"task-resource-attributes"},
			Short: taskResourceAttributesShort,
			Long:  taskResourceAttributesLong, PFlagProvider: taskresourceattribute.DefaultDelConfig, ProjectDomainNotRequired: true},
		"cluster-resource-attribute": {CmdFunc: deleteClusterResourceAttributes, Aliases: []string{"cluster-resource-attributes"},
			Short: clusterResourceAttributesShort,
			Long:  clusterResourceAttributesLong, PFlagProvider: clusterresourceattribute.DefaultDelConfig, ProjectDomainNotRequired: true},
		"execution-cluster-label": {CmdFunc: deleteExecutionClusterLabel, Aliases: []string{"execution-cluster-labels"},
			Short: executionClusterLabelShort,
			Long:  executionClusterLabelLong, PFlagProvider: executionclusterlabel.DefaultDelConfig, ProjectDomainNotRequired: true},
		"execution-queue-attribute": {CmdFunc: deleteExecutionQueueAttributes, Aliases: []string{"execution-queue-attributes"},
			Short: executionQueueAttributesShort,
			Long:  executionQueueAttributesLong, PFlagProvider: executionqueueattribute.DefaultDelConfig, ProjectDomainNotRequired: true},
	}
	cmdcore.AddCommands(deleteCmd, terminateResourcesFuncs)
	return deleteCmd
}
