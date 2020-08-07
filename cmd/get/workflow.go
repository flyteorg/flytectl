package get

import (
	"fmt"
	"context"
	"github.com/lyft/flytectl/cmd/config"
	"github.com/lyft/flytectl/pkg/printer"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytestdlib/logger"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
)

func getWorkflowFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	if config.GetConfig().Project == "" {
		return fmt.Errorf("Please set project name to get domain")
	}
	if config.GetConfig().Domain == "" {
		return fmt.Errorf("Please set project name to get workflow")
	}
	if len(args) > 0 {
		workflows, err := cmdCtx.AdminClient().ListWorkflows(ctx, &admin.ResourceListRequest{
			Id : &admin.NamedEntityIdentifier{
				Project: config.GetConfig().Project,
				Domain:  config.GetConfig().Domain,
				Name: args[0],
			},
			Limit:   3,
		})
		if err != nil {
			return err
		}
		logger.Debugf(ctx, "Retrieved %v workflows", len(workflows.Workflows))
		adminPrinter := printer.AdminWorkflowsList{
			Ctx: cmdCtx,
		}
		adminPrinter.Print(config.GetConfig().Output,workflows.Workflows)
		return nil
	}
	workflows, err := cmdCtx.AdminClient().ListWorkflowIds(ctx, &admin.NamedEntityIdentifierListRequest{
		Project: config.GetConfig().Project,
		Domain:  config.GetConfig().Domain,
		Limit: 3,
	})
	if err != nil {
		return err
	}
	logger.Debugf(ctx, "Retrieved %v workflows", len(workflows.Entities))
	adminPrinter := printer.AdminWorkflowsList{
		Ctx: cmdCtx,
	}
	adminPrinter.Print(config.GetConfig().Output,workflows.Entities)
	return nil
}

