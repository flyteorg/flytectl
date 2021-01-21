package delete

import (
	"context"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flytestdlib/logger"

	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
)


func terminateExecutionFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	name := args[0]
	logger.Infof(ctx, "Terminating execution of %v execution ", name)
	_, err := cmdCtx.AdminClient().TerminateExecution(ctx, &admin.ExecutionTerminateRequest{
		Id: &core.WorkflowExecutionIdentifier{
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
			Name: name,
		},
	})
	if err != nil {
		logger.Infof(ctx, "Failed in terminating execution of %v execution due to %v ", name, err)
		return err
	}
	logger.Infof(ctx, "Terminated execution of %v execution ", name)
	return nil
}
