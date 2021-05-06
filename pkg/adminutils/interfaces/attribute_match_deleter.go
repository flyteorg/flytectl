package interfaces

import (
	"context"

	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

type AttributeMatchDeleter interface {
	DeleteWorkflowAttributes(ctx context.Context, project, domain, name string, rsType admin.MatchableResource, cmdCtx cmdCore.CommandContext) (*admin.WorkflowAttributesDeleteResponse, error)
	DeleteProjectDomainAttributes(ctx context.Context, project, domain string, rsType admin.MatchableResource, cmdCtx cmdCore.CommandContext) (*admin.ProjectDomainAttributesDeleteResponse, error)
}

