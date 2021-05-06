package get

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
	taskResourceAttributesShort = "Gets matchable resources of task attributes"
	taskResourceAttributesLong  = `
Retrieves task  resource attributes for given project,domain combination or additionally with workflow name.

Usage
`
)

func getTaskResourceAttributes(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	var taskRsAttr *admin.TaskResourceAttributes
	if len(args) == 1 {
		workflowName := args[0]
		workflowAttr, err := cmdCtx.AdminFetcherExt().FetchWorkflowAttributes(ctx,
			config.GetConfig().Project, config.GetConfig().Domain, workflowName, admin.MatchableResource_TASK_RESOURCE)
		if err != nil {
			return err
		}
		taskRsAttr = workflowAttr.Attributes.MatchingAttributes.GetTaskResourceAttributes()
	} else {
		projectDomainAttr, err := cmdCtx.AdminFetcherExt().FetchProjectDomainAttributes(ctx,
			config.GetConfig().Project, config.GetConfig().Domain, admin.MatchableResource_TASK_RESOURCE)
		if err != nil {
			return err
		}
		taskRsAttr = projectDomainAttr.Attributes.MatchingAttributes.GetTaskResourceAttributes()
	}
	logger.Debugf(ctx, "Retrieved task attributes %v ", taskRsAttr)
	return nil
}
