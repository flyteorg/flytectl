package get

import (
	"context"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flytectl/pkg/printer"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flytestdlib/logger"
)

func getProjectsFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	if len(args) == 1 {
		projects, err := cmdCtx.AdminClient().ListProjects(ctx, &admin.ProjectListRequest{})
		if err != nil {
			return err
		}
		logger.Debugf(ctx, "Retrieved %v projects", len(projects.Projects))
		for _, v := range projects.Projects {
			if v.Name == args[0] {
				adminPrinter := printer.ProjectList{
					Ctx: cmdCtx,
				}
				adminPrinter.Print(config.GetConfig().Output, projects.Projects)
			}
		}
	}
	projects, err := cmdCtx.AdminClient().ListProjects(ctx, &admin.ProjectListRequest{})
	if err != nil {
		return err
	}
	logger.Debugf(ctx, "Retrieved %v projects", len(projects.Projects))
	adminPrinter := printer.ProjectList{
		Ctx: cmdCtx,
	}
	adminPrinter.Print(config.GetConfig().Output, projects.Projects)
	return nil
}
