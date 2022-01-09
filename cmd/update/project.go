package update

import (
	"context"
	"fmt"
	"github.com/flyteorg/flytectl/cmd/config/subcommand/project"
	"io/ioutil"

	"github.com/flyteorg/flytectl/clierrors"

	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"
)



const (
	projectShort = "Update project resources"
	projectLong  = `
Updates the project according to the flags passed. Allows you to archive or activate a project.
Activate project flytesnacks:
::

 flytectl update project -p flytesnacks --activateProject

Archive project flytesnacks:

::

 flytectl update project -p flytesnacks --archiveProject

Incorrect usage when passing both archive and activate:

::

 flytectl update project flytesnacks --archiveProject --activateProject

Incorrect usage when passing unknown-project:

::

 flytectl update project unknown-project --archiveProject

Incorrect usage when passing valid project using -p option:

::

 flytectl update project unknown-project --archiveProject -p known-project

Usage
`
)


func updateProjectsFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	id := config.GetConfig().Project
	projectSpec := project.CreateConfig{}
	if project.DefaultCreateConfig.File != "" {
		yamlFile, err := ioutil.ReadFile(project.DefaultUpdateConfig.File)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(yamlFile, &projectSpec)
		if err != nil {
			return err
		}
	} else {
		projectSpec.ID = id
		projectSpec.Description = project.DefaultUpdateConfig.Description
		projectSpec.Labels = project.DefaultUpdateConfig.Labels
	}
	if projectSpec.ID == "" {
		fmt.Errorf("project ID is required flag")
		return nil
	}

	archiveProject := project.DefaultUpdateConfig.ArchiveProject
	activateProject := project.DefaultUpdateConfig.ActivateProject
	if activateProject == archiveProject {
		return fmt.Errorf(clierrors.ErrInvalidStateUpdate)
	}
	projectState := admin.Project_ACTIVE
	if archiveProject {
		projectState = admin.Project_ARCHIVED
	}
	if project.DefaultUpdateConfig.DryRun {
		logger.Infof(ctx, "skipping UpdateProject request (dryRun)")
	} else {
		_, err := cmdCtx.AdminClient().UpdateProject(ctx, &admin.Project{
			Id:    projectSpec.ID,
			Description: projectSpec.Description,
			Labels: &admin.Labels{
				Values: projectSpec.Labels,
			},
			State: projectState,
		})
		if err != nil {
			fmt.Printf(clierrors.ErrFailedProjectUpdate, id, projectState, err)
			return err
		}
	}
	fmt.Printf("Project %v updated to %v state\n", id, projectState)
	return nil
}
