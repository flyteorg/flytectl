package create

import (
	"context"
	"fmt"
	"github.com/lyft/flytectl/cmd/config"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/lyft/flytestdlib/logger"

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


Usage
`
)

//go:generate pflags ProjectConfig

// ProjectConfig Config hold configuration for project create flags.
type ProjectConfig struct {
	ID          string `json:"id" pflag:",id of the project specified as argument."`
	Description string `json:"description" pflag:",description for the project specified as argument."`
}

var (
	projectConfig = &ProjectConfig{
		Description: "",
	}
)

func createProjectsCommand(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	id := config.GetConfig().Project
	if id == "" {
		fmt.Printf("Project not found")
		return nil
	}
	fmt.Printf("%v", projectConfig)

	response, err := cmdCtx.AdminClient().RegisterProject(ctx, &admin.ProjectRegisterRequest{
		Project: &admin.Project{
			Id:          projectConfig.ID,
			Name:        id,
			Description: projectConfig.Description,
		},
	})
	if err != nil {
		logger.Error(ctx, "Error %v", err)
	}
	logger.Debug(ctx, "Response %v", response)
	fmt.Println("Project Created successfully")
	return nil
}
