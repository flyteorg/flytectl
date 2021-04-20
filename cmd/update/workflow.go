package update

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/clierrors"
	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
)

const (
	updateWorkflowShort = "Updates launch plan metadata"
	updateWorkflowLong  = `
Updates workflow metadata.
Following command updates the description on the workflow.
::

 flytectl update workflow -p flytectldemo -d development core.advanced.run_merge_sort.merge_sort --description "Mergesort workflow example"

Archiving workflow named entity would cause this to disapper from flyteconsole UI.
::

 flytectl update workflow -p flytectldemo -d development  core.advanced.run_merge_sort.merge_sort --archive

Activating workflow named entity would unarchive it.
::

 flytectl update workflow -p flytectldemo -d development  core.advanced.run_merge_sort.merge_sort --activate

Usage
`
)

func updateWorkflowFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	project := config.GetConfig().Project
	domain := config.GetConfig().Domain
	if len(args) != 1 {
		return fmt.Errorf(clierrors.ErrWorkflowNotPassed)
	}
	name := args[0]
	err := namedEntityConfig.UpdateNamedEntity(ctx, name, project, domain, core.ResourceType_WORKFLOW, cmdCtx)
	if err != nil {
		fmt.Printf(clierrors.ErrFailedWorkflowUpdate, name, err)
		return err
	}
	fmt.Printf("updated metadata successfully on %v", name)
	return nil
}
