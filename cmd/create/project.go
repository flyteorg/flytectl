package create

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/flyteorg/flytectl/cmd/config/subcommand/project"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"

	"gopkg.in/yaml.v2"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytestdlib/logger"
)

const (
	projectShort = "Create project resources"
	projectLong  = `
Create projects.(project/projects can be used interchangeably in these commands)

::

 flytectl create project --name flytesnacks --id flytesnacks --description "flytesnacks description"  --labels app=flyte

Create a project by definition file. Note: The name shouldn't contain any whitespace characters.
::

 flytectl create project --file project.yaml 

.. code-block:: yaml

    id: "project-unique-id"
    name: "Name"
    labels:
     app: flyte
    description: "Some description for the project"

`
)

func createProjectsCommand(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	projectSpec := project.Definition{}
	if project.DefaultProjectConfig.File != "" {
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
		projectSpec.Name = project.DefaultProjectConfig.Name
		projectSpec.Description = project.DefaultProjectConfig.Description
		projectSpec.Labels = project.DefaultProjectConfig.Labels
	}
	if projectSpec.ID == "" {
		return fmt.Errorf("project ID is required flag")
	}
	if projectSpec.Name == "" {
		return fmt.Errorf("project name is required flag")
	}

	if project.DefaultProjectConfig.DryRun {
		logger.Debugf(ctx, "skipping RegisterProject request (DryRun)")
	} else {
		_, err := cmdCtx.AdminClient().RegisterProject(ctx, &admin.ProjectRegisterRequest{
			Project: &admin.Project{
				Id:          projectSpec.ID,
				Name:        projectSpec.Name,
				Description: projectSpec.Description,
				Labels: &admin.Labels{
					Values: projectSpec.Labels,
				},
			},
		})
		if err != nil {
			return err
		}
	}
	fmt.Println("project Created successfully")
	return nil
}
