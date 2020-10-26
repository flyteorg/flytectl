package create

import (
	"context"
	"fmt"
	"github.com/lyft/flytectl/cmd/config"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flytestdlib/logger"
	"io/ioutil"
	"gopkg.in/yaml.v2"


	cmdCore "github.com/lyft/flytectl/cmd/core"
)

func createProjectsFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	var project CreateProject
	if config.GetCreateConfig().Filename != "" {
		yamlFile, err := ioutil.ReadFile(config.GetCreateConfig().Filename)
		if err != nil {
			logger.Error(ctx, "Error %v", err)
		}
		err = yaml.Unmarshal(yamlFile, project)
		if err != nil {
			logger.Error(ctx, "Error %v", err)
		}
	}else{
		project.Name = config.GetCreateConfig().Name
		if project.Name == "" {
			logger.Debug(ctx, "Name is required to create a project")
		}
		project.Id = config.GetCreateConfig().ID
		if project.Id == "" {
			logger.Debug(ctx, "Id is required to create a project")
		}
		project.Labels = config.GetCreateConfig().Labels
		project.Description = config.GetCreateConfig().Description
	}
	response,err := cmdCtx.AdminClient().RegisterProject(ctx,&admin.ProjectRegisterRequest{
		Project: &admin.Project{
			Id: project.Id,
			Name: project.Name,
			Description: project.Description,
		},
	},)
	if err != nil {
		logger.Error(ctx, "Error %v", err)
	}
	logger.Debug(ctx,"Response %v",response)
	fmt.Println("OK")
	return nil
}
