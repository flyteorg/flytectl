package update

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/clierrors"
	"github.com/flyteorg/flytectl/cmd/config/subcommand/project"
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

 flytectl update project -p flytesnacks --activate

Archive project flytesnacks:

::

 flytectl update project -p flytesnacks --archive

Incorrect usage when passing both archive and activate:

::

 flytectl update project -p flytesnacks --archiveProject --activate

Incorrect usage when passing unknown-project:

::

 flytectl update project unknown-project --archive

project ID is required flag

::

 flytectl update project unknown-project --archiveProject

Update projects.(project/projects can be used interchangeably in these commands)

::

 flytectl update project -p flytesnacks --description "flytesnacks description"  --labels app=flyte

Update a project by definition file. Note: The name shouldn't contain any whitespace characters.
::

 flytectl update project --file project.yaml 

.. code-block:: yaml

    id: "project-unique-id"
    name: "Name"
    labels:
       values:
         app: flyte
    description: "Some description for the project"

Update a project state by definition file. Note: The name shouldn't contain any whitespace characters.
::

 flytectl update project --file project.yaml  --archive

.. code-block:: yaml

    id: "project-unique-id"
    name: "Name"
    labels:
       values:
         app: flyte
    description: "Some description for the project"

Usage
`
)

func updateProjectsFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	projectSpec, err := project.DefaultProjectConfig.GetProjectSpec(project.DefaultProjectConfig.ID)
	if err != nil {
		return err
	}
	if projectSpec.Id == "" {
		return fmt.Errorf(clierrors.ErrProjectNotPassed)
	}

	projectDefinition := &admin.Project{
		Id: projectSpec.Id,
	}
	if projectSpec.Description != "" {
		projectDefinition.Description = projectSpec.Description
	}
	if projectSpec.Name != "" {
		projectDefinition.Name = projectSpec.Name
	}
	if len(projectSpec.Labels.Values) > 0 {
		projectDefinition.Labels = projectSpec.Labels
	}

	projectDefinition, err = project.DefaultProjectConfig.MapToAdminState(projectDefinition)
	if err != nil {
		return err
	}

	if project.DefaultProjectConfig.DryRun {
		logger.Infof(ctx, "skipping UpdateProject request (dryRun)")
	} else {
		_, err := cmdCtx.AdminClient().UpdateProject(ctx, projectDefinition)
		if err != nil {
			fmt.Printf(clierrors.ErrFailedProjectUpdate, projectSpec.Id, err)
			return err
		}
	}
	fmt.Printf("Project %v updated\n", projectSpec.Id)
	return nil
}
