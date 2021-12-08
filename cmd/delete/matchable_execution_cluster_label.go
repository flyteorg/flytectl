package delete

import (
	"context"

	"github.com/flyteorg/flytectl/cmd/config"
	sconfig "github.com/flyteorg/flytectl/cmd/config/subcommand"
	"github.com/flyteorg/flytectl/cmd/config/subcommand/executionclusterlabel"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

const (
	executionClusterLabelShort = "Deletes matchable resources of execution cluster label"
	executionClusterLabelLong  = `
Deletes execution cluster label for given project and domain combination or additionally with workflow name.

Deletes execution cluster label for project and domain.
The command delete execution cluster label for project FlyteCTLdemo and development domain:
::

 flytectl delete execution-cluster-label -p flytectldemo -d development 


Deletes execution cluster label using config file which was used for creating it.
The command deletes execution cluster label from the config file ecl.yaml.
Value is optional in the file as it is unread during the delete command but it can be kept since the same file can be used for get, update or delete commands. 
Ex:  Content of ecl.yaml which will use the project domain and workflow name to delete the resource.

::

 flytectl delete execution-cluster-label --attrFile ecl.yaml


.. code-block:: yaml
	
    domain: development
    project: flytectldemo
    value: foo

Deletes execution cluster label for a workflow.
The command deletes execution cluster label for a workflow core.control_flow.run_merge_sort.merge_sort:

::

 flytectl delete execution-cluster-label -p flytectldemo -d development core.control_flow.run_merge_sort.merge_sort

Usage
`
)

func deleteExecutionClusterLabel(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	var pwdGetter sconfig.ProjectDomainWorkflowGetter
	pwdGetter = sconfig.PDWGetterCommandLine{Config: config.GetConfig(), Args: args}
	delConfig := executionclusterlabel.DefaultDelConfig

	// Get the project domain workflowName from the config file or commandline params
	if len(delConfig.AttrFile) > 0 {
		// Initialize FileConfig which will be used if delConfig.AttrFile is non empty
		// And Reads from the cluster label file
		pwdGetter = &executionclusterlabel.FileConfig{}
		if err := sconfig.ReadConfigFromFile(pwdGetter, delConfig.AttrFile); err != nil {
			return err
		}
	}
	// Use the pwdGetter to initialize the project domain and workflow
	project := pwdGetter.GetProject()
	domain := pwdGetter.GetDomain()
	workflowName := pwdGetter.GetWorkflow()

	// Deletes the matchable attributes using the ExecClusterLabelFileConfig
	if err := deleteMatchableAttr(ctx, project, domain, workflowName, cmdCtx.AdminDeleterExt(),
		admin.MatchableResource_EXECUTION_CLUSTER_LABEL, delConfig.DryRun); err != nil {
		return err
	}

	return nil
}
