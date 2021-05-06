package delete

import (
	"context"
	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"
)

//go:generate pflags MatchableRsAttrConfig --default-var matchableRsConfig
var (
	matchableRsConfig = &TaskResourceAttrConfig{}
)

// TaskResourceAttrConfig Matchable resource attributes configuration passed from command line
type TaskResourceAttrConfig struct {
}

const (
	taskResourceAttributesShort = "Deletes matchable resources of task attributes"
	taskResourceAttributesLong  = `
Deletes task  resource attributes for given project,domain combination or additionally with workflow name.

Usage
`
)

func deleteTaskResourceAttributes(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	project := config.GetConfig().Project
	domain := config.GetConfig().Domain
	if len(args) == 1 {
		workflowName := args[0]
		err := cmdCtx.AdminDeleterExt().DeleteWorkflowAttributes(ctx, project, domain, workflowName, admin.MatchableResource_TASK_RESOURCE)
		if err != nil {
			return err
		}
		logger.Debugf(ctx, "Deleted task resource attributes from %v project and domain %v and workflow %v", project, domain, workflowName)
	} else {
		err := cmdCtx.AdminDeleterExt().DeleteProjectDomainAttributes(ctx, project, domain, admin.MatchableResource_TASK_RESOURCE)
		if err != nil {
			return err
		}

		logger.Debugf(ctx, "Deleted task resource attributes from %v project and domain %v", project, domain)
	}
	return nil
}
