package update

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/clierrors"
	"github.com/flyteorg/flytectl/cmd/config"
	"github.com/flyteorg/flytectl/cmd/config/subcommand/execution"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/flyteorg/flytestdlib/logger"
)

const (
	updateExecutionShort = "Update execution status"
	updateExecutionLong  = `
Activating execution unhides it from cli and UI:
::

 flytectl update execution -p flytectldemo -d development  oeh94k9r2r --activate

Archiving execution hides it from cli and UI:
::

 flytectl update execution -p flytectldemo -d development  oeh94k9r2r --archive


Usage
`
)

func updateExecutionFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	project := config.GetConfig().Project
	domain := config.GetConfig().Domain
	if len(args) != 1 {
		return fmt.Errorf(clierrors.ErrExecutionNotPassed)
	}
	executionName := args[0]
	activateLP := execution.UConfig.Activate
	archiveLP := execution.UConfig.Archive
	if activateLP == archiveLP && archiveLP {
		return fmt.Errorf(clierrors.ErrInvalidStateUpdate)
	}

	var executionState admin.ExecutionStatus_ExecutionState
	if activateLP {
		executionState = admin.ExecutionStatus_EXECUTION_ACTIVE
	} else if archiveLP {
		executionState = admin.ExecutionStatus_EXECUTION_ARCHIVED
	}

	if execution.UConfig.DryRun {
		logger.Debugf(ctx, "skipping UpdateExecution request (DryRun)")
	} else {
		_, err := cmdCtx.AdminClient().UpdateExecution(ctx, &admin.ExecutionUpdateRequest{
			Id: &core.WorkflowExecutionIdentifier{
				Project: project,
				Domain:  domain,
				Name:    executionName,
			},
			Status: &admin.ExecutionStatus{State: executionState},
		})
		if err != nil {
			fmt.Printf(clierrors.ErrFailedExecutionUpdate, executionName, err)
			return err
		}
	}
	fmt.Printf("updated execution successfully on %v to state %v", executionName, executionState)

	return nil
}
