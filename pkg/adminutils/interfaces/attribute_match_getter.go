package interfaces

import (
	"context"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

type AttributeMatchGetter interface {
	GetWorkflowAttributes(ctx context.Context, adminClient service.AdminServiceClient, project, domain, name string,
		rsType admin.MatchableResource) (*admin.WorkflowAttributesGetResponse, error)
	GetProjectDomainAttributes(ctx context.Context, adminClient service.AdminServiceClient, project, domain string,
		rsType admin.MatchableResource) (*admin.ProjectDomainAttributesGetResponse, error)
}

