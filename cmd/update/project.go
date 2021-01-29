package update

import (
	"context"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flytestdlib/logger"
)

func archiveProjectFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	return updateProjectsFunc(ctx, admin.Project_ARCHIVED, cmdCtx)
}

func activateProjectFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	return updateProjectsFunc(ctx, admin.Project_ACTIVE, cmdCtx)
}

func updateProjectsFunc(ctx context.Context, projectState admin.Project_ProjectState, cmdCtx cmdCore.CommandContext) error {
	id := config.GetConfig().Project
	_, err := cmdCtx.AdminClient().UpdateProject(ctx, &admin.Project{
		Id : id,
		State : projectState,
	})
	if err != nil {
		return err
	}
	logger.Infof(ctx, "Project %v updated to %v state", id, projectState)
	return nil
}
