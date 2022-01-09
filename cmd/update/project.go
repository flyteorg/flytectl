package update

import (
	"context"
	"fmt"

	"github.com/flyteorg/flytectl/clierrors"
	"github.com/flyteorg/flytectl/cmd/config/subcommand/project"
	"gopkg.in/yaml.v2"

	"io/ioutil"

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

 flytectl update project --id flytesnacks --activateProject

Archive project flytesnacks:

::

 flytectl update project --id flytesnacks --archiveProject

Incorrect usage when passing both archive and activate:

::

 flytectl update project flytesnacks --archiveProject --activateProject

Incorrect usage when passing unknown-project:

::

 flytectl update project unknown-project --archiveProject

project ID is required flag

::

 flytectl update project unknown-project --archiveProject -p known-project

Update projects.(project/projects can be used interchangeably in these commands)

::

 flytectl update project --id flytesnacks --description "flytesnacks description"  --labels app=flyte

Update a project by definition file. Note: The name shouldn't contain any whitespace characters.
::

 flytectl update project --file project.yaml 

.. code-block:: yaml

    id: "project-unique-id"
    name: "Name"
    labels:
     app: flyte
    description: "Some description for the project"

Usage
`
)

func updateProjectsFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	projectSpec := project.Definition{}
	if len(project.DefaultProjectConfig.File) > 0 {
		yamlFile, err := ioutil.ReadFile(project.DefaultProjectConfig.File)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(yamlFile, &projectSpec)
		if err != nil {
			return err
		}
	} else {
		projectSpec.ID = project.DefaultProjectConfig.ID
		projectSpec.Description = project.DefaultProjectConfig.Description
		projectSpec.Name = project.DefaultProjectConfig.Name
		projectSpec.Labels = project.DefaultProjectConfig.Labels
	}

	if projectSpec.ID == "" {
		return fmt.Errorf("project ID is required flag")
	}

	projectDefinition := &admin.Project{
		Id: projectSpec.ID,
	}
	if projectSpec.Description != "" {
		projectDefinition.Description = projectSpec.Description
	}
	if projectSpec.Name != "" {
		projectDefinition.Name = projectSpec.Name
	}
	if len(projectSpec.Labels) > 0 {
		projectDefinition.Labels = &admin.Labels{
			Values: projectSpec.Labels,
		}
	}

	activateProject := project.DefaultProjectConfig.ActivateProject
	archiveProject := project.DefaultProjectConfig.ArchiveProject
	if activateProject || archiveProject {
		if activateProject == archiveProject {
			return fmt.Errorf(clierrors.ErrInvalidStateUpdate)
		}
		projectDefinition.State = admin.Project_ACTIVE
		if archiveProject {
			projectDefinition.State = admin.Project_ARCHIVED
		}
	}

	if project.DefaultProjectConfig.DryRun {
		logger.Infof(ctx, "skipping UpdateProject request (dryRun)")
	} else {
		_, err := cmdCtx.AdminClient().UpdateProject(ctx, projectDefinition)
		if err != nil {
			fmt.Printf(clierrors.ErrFailedProjectUpdate, projectSpec.ID, err)
			return err
		}
	}
	fmt.Printf("Project %v updated\n", projectSpec.ID)
	return nil
}
