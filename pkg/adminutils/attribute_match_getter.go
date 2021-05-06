package adminutils

import (
	"context"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

type AttributeMatchGetterConfig struct {
	AttrFile string `json:"attrFile" pflag:",attribute file name to be used for generating attribute for the resource type."`
	ResourceType  string `json:"resourceType" pflag:", type of resource for which attributes need to be fetched.It can take following values TASK_RESOURCE CLUSTER_RESOURCE EXECUTION_QUEUE EXECUTION_CLUSTER_LABEL QUALITY_OF_SERVICE_SPECIFICATION PLUGIN_OVERRIDE"`
}

func (am AttributeMatchGetterConfig) GetWorkflowAttributes(ctx context.Context, project, domain, name string, rsType admin.MatchableResource, cmdCtx cmdCore.CommandContext) (*admin.WorkflowAttributesGetResponse, error) {
	workflowAttr, err := cmdCtx.AdminClient().GetWorkflowAttributes(ctx, &admin.WorkflowAttributesGetRequest{
		Project: project,
		Domain: domain,
		Workflow: name,
		ResourceType: rsType,
	})
	return workflowAttr, err
}

func (am AttributeMatchGetterConfig) GetProjectDomainAttributes(ctx context.Context, project, domain string, rsType admin.MatchableResource, cmdCtx cmdCore.CommandContext) (*admin.ProjectDomainAttributesGetResponse, error) {
	projectDomainAttr, err := cmdCtx.AdminClient().GetProjectDomainAttributes(ctx, &admin.ProjectDomainAttributesGetRequest{
		Project: project,
		Domain: domain,
		ResourceType: rsType,
	})
	return projectDomainAttr, err
}

//
//func (am AttributeMatchGetterConfig) UpdateWorkflowAttributes(ctx context.Context, project, domain, name string, rsType admin.MatchableResource, cmdCtx cmdCore.CommandContext) (*admin.WorkflowAttributesUpdateResponse, error) {
//	workflowAttr, err := cmdCtx.AdminClient().UpdateWorkflowAttributes(ctx, &admin.WorkflowAttributesUpdateRequest{
//		Attributes: &admin.WorkflowAttributes{
//			Project: project,
//			Domain: domain,
//			Workflow: name,
//			MatchingAttributes: &admin.MatchingAttributes{
//				Target: &admin.MatchingAttributes_TaskResourceAttributes{
//					TaskResourceAttributes: &admin.TaskResourceAttributes{
//						Defaults: &admin.TaskResourceSpec{
//							Cpu: "1m",
//							Gpu: "1m",
//							Memory: "1Mi",
//						},
//					},
//				},
//			},
//		},
//	})
//	return workflowAttr, err
//}
//
//func (am AttributeMatchGetterConfig) UpdateProjectDomainAttributes(ctx context.Context, project, domain string, rsType admin.MatchableResource, cmdCtx cmdCore.CommandContext) (*admin.ProjectDomainAttributesUpdateResponse, error) {
//	projectDomainAttr, err := cmdCtx.AdminClient().UpdateProjectDomainAttributes(ctx, &admin.ProjectDomainAttributesUpdateRequest{
//		Attributes: &admin.ProjectDomainAttributes{
//			Project: project,
//			Domain: domain,
//			MatchingAttributes: &admin.MatchingAttributes{
//				Target: &admin.MatchingAttributes_TaskResourceAttributes{
//					TaskResourceAttributes: &admin.TaskResourceAttributes{
//						Defaults: &admin.TaskResourceSpec{
//							Cpu: "1m",
//							Gpu: "1m",
//							Memory: "1Mi",
//						},
//					},
//				},
//			},
//		},
//	})
//	return projectDomainAttr, err
//}