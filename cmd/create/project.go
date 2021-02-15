package create

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flytestdlib/logger"
	"gopkg.in/yaml.v2"

	cmdCore "github.com/lyft/flytectl/cmd/core"
)

const (
	projectShort = "Create project resources"
	projectLong  = `
Create the projects.(project,projects can be used interchangeably in these commands)
::

 bin/flytectl create project --id test --description test -p test
Project Created

::

Create the project using yaml definition file

::
 bin/flytectl create project --file project.yaml 
Project Created successfully

::

Usage
`
)

//go:generate pflags ProjectConfig --default-var projectConfig --bind-default-var

// ProjectConfig Config hold configuration for project create flags.
type ProjectConfig struct {
	id          string `json:"id" pflag:",id for the project specified as argument."`
	name        string `json:"name" pflag:",name for the project specified as argument."`
	file        string `json:"file" pflag:",file for the project definition."`
	description string `json:"description" pflag:",description for the project specified as argument."`
	labels map[string]string `json:"labels" pflag:",labels for the project specified as argument."`
}

var (
	projectConfig = &ProjectConfig{
		description: "",
	}
)

func createProjectsCommand(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	project := projectDefinition{}
	if projectConfig.file != "" {
		yamlFile, err := ioutil.ReadFile(projectConfig.file)
		if err != nil {
			logger.Error(ctx, "Error %v", err)
		}
		err = yaml.Unmarshal(yamlFile, &project)
		if err != nil {
			logger.Error(ctx, "Error %v", err)
		}
	} else {
		project.ID = projectConfig.id
		project.Name = projectConfig.name
		project.Description = projectConfig.description
	}
	if project.ID == "" {
		fmt.Printf("Project ID is required flag")
		return nil
	}
	if project.Name == "" {
		fmt.Printf("Project name is required flag")
		return nil
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
		fmt.Printf("error: %v", err.Error())
		return nil
	}
	fmt.Println("Project Created successfully")
	return nil
}
