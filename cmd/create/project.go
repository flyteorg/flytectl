package create

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"gopkg.in/yaml.v2"

	cmdCore "github.com/lyft/flytectl/cmd/core"
)

const (
	projectShort = "Create project resources"
	projectLong  = `
Create the projects.(project,projects can be used interchangeably in these commands)
::

 bin/flytectl create project --name flytesnacks --id flytesnacks --description "flytesnacks description"  --labels app=flyte
Project Created

::

Create the project using yaml definition file

::
 bin/flytectl create project --file project.yaml 
Project Created successfully

::

.. code-block:: yaml

   id: "project-unique-id"
   name: "Friendly name"
   labels:
	  app: flyte
   description: "Some description for the project"
Usage
`
)

//go:generate pflags ProjectConfig --default-var projectConfig --bind-default-var

// ProjectConfig Config hold configuration for project create flags.
type ProjectConfig struct {
	ID          string            `json:"id" pflag:",id for the project specified as argument."`
	Name        string            `json:"name" pflag:",name for the project specified as argument."`
	File        string            `json:"file" pflag:",file for the project definition."`
	Description string            `json:"description" pflag:",description for the project specified as argument."`
	Labels      map[string]string `json:"labels" pflag:",labels for the project specified as argument."`
}

var (
	projectConfig = &ProjectConfig{
		Description: "",
		Labels:      map[string]string{},
	}
)

func createProjectsCommand(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	project := projectDefinition{}
	if projectConfig.File != "" {
		yamlFile, err := ioutil.ReadFile(projectConfig.File)
		if err != nil {
			return fmt.Errorf("Error %v", err)
		}
		err = yaml.Unmarshal(yamlFile, &project)
		if err != nil {
			return fmt.Errorf("Error %v", err)
		}
	} else {
		project.ID = projectConfig.ID
		project.Name = projectConfig.Name
		project.Description = projectConfig.Description
		project.Labels = projectConfig.Labels
	}
	if project.ID == "" {
		fmt.Printf("project ID is required flag")
		return fmt.Errorf("project ID is required flag")
	}
	if project.Name == "" {
		return fmt.Errorf("project name is required flag")
	}
	_, err := cmdCtx.AdminClient().RegisterProject(ctx, &admin.ProjectRegisterRequest{
		Project: &admin.Project{
			Id:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			Labels: &admin.Labels{
				Values: project.Labels,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error: %v", err.Error())
	}
	fmt.Println("project Created successfully")
	return nil
}
