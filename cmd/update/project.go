package update

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/cmd/config"
	"github.com/flyteorg/flytectl/pkg/util"

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

 flytectl update project --id flytesnacks --activate

Archive project flytesnacks:

::

 flytectl update project --id flytesnacks --archive

Incorrect usage when passing both archive and activate:

::

 flytectl update project flytesnacks --archiveProject --activate

Incorrect usage when passing unknown-project:

::

 flytectl update project unknown-project --archive

project ID is required flag

::

 flytectl update project unknown-project --archiveProject -p known-project

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
	projectSpec, err := util.GetProjectSpec(project.DefaultProjectConfig, config.GetConfig().Project)
	if err != nil {
		return err
	}
	if projectSpec.Id == "" {
		return fmt.Errorf("project ID is required flag")
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
	if projectSpec.Labels != nil {
		projectDefinition.Labels = projectSpec.Labels
	}

	projectDefinition, err = getState(project.DefaultProjectConfig, projectDefinition)
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

func getState(flags *project.ConfigProject, spec *admin.Project) (*admin.Project, error) {
	if flags.ActivateProject {
		flags.Activate = flags.ActivateProject
	}
	if flags.ArchiveProject {
		flags.Archive = flags.ArchiveProject
	}

	activate := flags.Activate
	archive := flags.Archive

	if activate || archive {
		if activate == archive {
			return spec, fmt.Errorf(clierrors.ErrInvalidStateUpdate)
		}
		spec.State = admin.Project_ACTIVE
		if activate {
			spec.State = admin.Project_ARCHIVED
		}
		return spec, nil
	}
	return spec, nil
}
