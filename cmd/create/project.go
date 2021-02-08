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

//go:generate pflags ProjectConfig

// Config hold configuration for project create flags.
type ProjectConfig struct {
	Name        string        `json:"name" pflag:",Name of the project specified as argument."`
	ID          string        `json:"id" pflag:",Id of the project specified as argument."`
	Filename    string        `json:"file" pflag:",Filename of the project specified as argument."`
	Labels      *admin.Labels `json:"labels" pflag:",Labels for the project specified as argument."`
	Description string        `json:"description" pflag:",Description for the project specified as argument."`
}

var (
	projectConfig = &ProjectConfig{}
)

func createProjectsFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {

	pconfig := ProjectConfig{}

	if projectConfig.Filename != "" {
		yamlFile, err := ioutil.ReadFile(projectConfig.Filename)
		if err != nil {
			logger.Error(ctx, "Error %v", err)
		}
		err = yaml.Unmarshal(yamlFile, pconfig)
		if err != nil {
			logger.Error(ctx, "Error %v", err)
		}
	} else {
		pconfig.Name = projectConfig.Name
		if pconfig.Name == "" {
			logger.Debug(ctx, "Name is required to create a project")
		}
		pconfig.ID = projectConfig.ID
		if pconfig.ID == "" {
			logger.Debug(ctx, "Id is required to create a project")
		}
		pconfig.Labels = projectConfig.Labels
		pconfig.Description = projectConfig.Description
	}
	response, err := cmdCtx.AdminClient().RegisterProject(ctx, &admin.ProjectRegisterRequest{
		Project: &admin.Project{
			Id:          pconfig.ID,
			Name:        pconfig.Name,
			Description: pconfig.Description,
			Labels:      pconfig.Labels,
		},
	})
	if err != nil {
		logger.Error(ctx, "Error %v", err)
	}
	logger.Debug(ctx, "Response %v", response)
	fmt.Println("Project Created successfully")
	return nil
}
