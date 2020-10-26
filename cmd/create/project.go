package create

import (
	"context"
	"github.com/lyft/flytectl/cmd/config"
	"github.com/lyft/flytestdlib/logger"
	"io/ioutil"

	cmdCore "github.com/lyft/flytectl/cmd/core"
)

func createProjectsFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	if config.GetConfig().Filename != "" {
		yamlFile, err := ioutil.ReadFile("conf.yaml")
		if err != nil {
			logger.Error(ctx, "Error %v", err)
		}
		err = yaml.Unmarshal(yamlFile, projects)
	}
	return nil
}
