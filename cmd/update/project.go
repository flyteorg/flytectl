package update

import (
	"context"
	"fmt"
	"os"

	"github.com/flyteorg/flytectl/clierrors"
	"github.com/flyteorg/flytectl/cmd/config"
	"github.com/flyteorg/flytectl/cmd/config/subcommand/project"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	cmdUtil "github.com/flyteorg/flytectl/pkg/commandutils"
)

const (
	projectShort = "Update project resources"
	projectLong  = `
Update the project according to the flags passed. Allows you to archive or activate a project.
Activate project flytesnacks:
::

 flytectl update project -p flytesnacks --activate

Archive project flytesnacks:

::

 flytectl update project -p flytesnacks --archive

Incorrect usage when passing both archive and activate:

::

 flytectl update project -p flytesnacks --archive --activate

Incorrect usage when passing unknown-project:

::

 flytectl update project unknown-project --archive

project ID is required flag

::

 flytectl update project unknown-project --archive

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
	newProject, err := project.DefaultProjectConfig.GetProjectSpec(config.GetConfig())
	if err != nil {
		return err
	}

	if newProject.Id == "" {
		return fmt.Errorf(clierrors.ErrProjectNotPassed)
	}

	oldProject, err := cmdCtx.AdminFetcherExt().GetProjectById(ctx, newProject.Id)
	if err != nil {
		fmt.Printf(clierrors.ErrFailedProjectUpdate, newProject.Id, err)
		return err
	}

	patch, err := diffAsYaml(oldProject, newProject)
	if err != nil {
		panic(err)
	}

	if patch == "" {
		fmt.Printf("No changes detected. Skipping the update.\n")
		return nil
	}

	fmt.Printf("The following changes are to be applied.\n%s\n", patch)

	if project.DefaultProjectConfig.DryRun {
		fmt.Printf("skipping UpdateProject request (dryRun)\n")
		return nil
	}

	if !project.DefaultProjectConfig.Force && !cmdUtil.AskForConfirmation("Continue?", os.Stdin) {
		return fmt.Errorf("update aborted")
	}

	_, err = cmdCtx.AdminClient().UpdateProject(ctx, newProject)
	if err != nil {
		fmt.Printf(clierrors.ErrFailedProjectUpdate, newProject.Id, err)
		return err
	}

	fmt.Printf("project %s updated\n", newProject.Id)
	return nil
}
