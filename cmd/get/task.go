package get

import (
	"context"
	"fmt"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/pkg/printer"
	"github.com/lyft/flytestdlib/logger"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
)

func getTaskFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	if config.GetConfig().Project == "" {
		return fmt.Errorf("Please set project name to get domain")
	}
	if config.GetConfig().Domain == "" {
		return fmt.Errorf("Please set project name to get workflow")
	}
	if len(args) == 1 {
		task, err := cmdCtx.AdminClient().ListTasks(ctx, &admin.ResourceListRequest{
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
		logger.Debugf(ctx, "Retrieved Task",task.Tasks)
		taskPrinter := printer.AdminTasksList{
			Ctx: cmdCtx,
		}
		taskPrinter.Print(config.GetConfig().Output,task.Tasks)
		return nil
	}

	tasks, err := cmdCtx.AdminClient().ListTaskIds(ctx, &admin.NamedEntityIdentifierListRequest{
		Project: config.GetConfig().Project,
		Domain:  config.GetConfig().Domain,
		Limit: 10,
	})
	if err != nil {
		return err
	}
	logger.Debugf(ctx, "Retrieved %v Task", len(tasks.Entities))
	taskPrinter := printer.AdminTasksList{
		Ctx: cmdCtx,
	}
	taskPrinter.Print(config.GetConfig().Output,tasks.Entities)
	return nil
}
