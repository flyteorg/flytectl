package update

import (
	"context"
	"github.com/flyteorg/flytectl/cmd/config"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flytestdlib/logger"
)

//go:generate pflags TaskResourceAttrConfig --default-var matchableTaskRsConfig
var (
	matchableTaskRsConfig = &TaskResourceAttrConfig{}
)

// TaskResourceAttrConfig Matchable resource attributes configuration passed from command line
type TaskResourceAttrConfig struct {
	AttrFile       string `json:"attrFile" pflag:",attribute file name to be used for generating attribute for the resource type."`
	CpuDefault     string `json:"cpuDefault" pflag:",cpuDefault to be assigned for workflow or tasks."`
	GpuDefault     string `json:"gpuDefault" pflag:",gpuDefault to be assigned for workflow or tasks."`
	MemoryDefault  string `json:"memoryDefault" pflag:",memoryDefault to be assigned for workflow or tasks."`
	StorageDefault string `json:"storageDefault" pflag:",storageDefault to be assigned for workflow or tasks."`
	CpuLimit       string `json:"cpuLimit" pflag:",cpuLimit to be assigned for workflow or tasks."`
	GpuLimit       string `json:"gpuLimit" pflag:",gpuLimit to be assigned for workflow or tasks."`
	MemoryLimit    string `json:"memoryLimit" pflag:",memoryLimit to be assigned for workflow or tasks."`
	StorageLimit   string `json:"storageLimit" pflag:",storageLimit to be assigned for workflow or tasks."`
}

const (
	taskResourceAttributesShort = "Updates matchable resources of task attributes"
	taskResourceAttributesLong  = `
Updates task  resource attributes for given project,domain combination or additionally with workflow name.

Usage
`
)

func updateTaskResourceAttributesFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	project := config.GetConfig().Project
	domain := config.GetConfig().Domain
	if len(args) == 1 {
		workflowName := args[0]
		err := cmdCtx.AdminUpdaterExt().UpdateWorkflowAttributes(ctx, project, domain, workflowName, admin.MatchableResource_TASK_RESOURCE)
		if err != nil {
			return err
		}

		logger.Debugf(ctx, "Updated task resource attributes from %v project and domain %v and workflow %v", project, domain, workflowName)
	} else {
		err := cmdCtx.AdminUpdaterExt().UpdateProjectDomainAttributes(ctx, project, domain, admin.MatchableResource_TASK_RESOURCE)
		if err != nil {
			return err
		}

		logger.Debugf(ctx, "Updated task resource attributes from %v project and domain %v", project, domain)
	}
	return nil
}
