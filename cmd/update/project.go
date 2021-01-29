package update

import (
	"context"
	"errors"
	"github.com/lyft/flytectl/cmd/config"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flytestdlib/logger"
)

var (
	errProjectNotFound = errors.New("Specify id of the project to be updated")
	errInvalidUpdate = errors.New("Specify either activate or archive")
)

func updateProjectsFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	id := config.GetConfig().Project
	if len(id) == 0 {
		return errProjectNotFound
	}
	archiveProject := GetConfig().ArchiveProject
	activateProject := GetConfig().ActivateProject
	if activateProject == archiveProject {
		return errInvalidUpdate
	}
	projectState := admin.Project_ACTIVE
	if archiveProject {
		projectState = admin.Project_ARCHIVED
	}
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
