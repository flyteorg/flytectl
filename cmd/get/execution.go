package get

import (
	"context"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/pkg/printer"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flytestdlib/logger"
)

func getExecutionFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	executionPrinter := printer.Printer{}
	excutions, err := cmdCtx.AdminClient().ListExecutions(ctx, &admin.ResourceListRequest{
		Limit: 10,
		Id: &admin.NamedEntityIdentifier{
			Project: config.GetConfig().Project,
			Domain:  config.GetConfig().Domain,
		},
	})
	if err != nil {
		return err
	}
	if len(args) == 1 {
		name := args[0]
		if err != nil {
			return err
		}
		for _, v := range excutions.Executions {
			if v.Id.Name == name {
				err := executionPrinter.Print(config.GetConfig().MustOutputFormat(), v, executionSingleStructure, transformSingleExecution)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return nil
	}

	logger.Debugf(ctx, "Retrieved %v excutions", len(excutions.Executions))
	executionPrinter.Print(config.GetConfig().MustOutputFormat(), excutions.Executions, executionStructure, transformExecution)
	return nil
}
