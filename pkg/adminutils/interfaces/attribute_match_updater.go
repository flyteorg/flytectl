package interfaces

import (
	"context"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

type AttributeMatchUpdater interface {
	UpdateWorkflowAttributes(ctx context.Context, adminClient service.AdminServiceClient, project, domain, name string,
		rsType admin.MatchableResource) (*admin.WorkflowAttributesUpdateResponse, error)
	UpdateProjectDomainAttributes(ctx context.Context, adminClient service.AdminServiceClient, project, domain string,
		rsType admin.MatchableResource) (*admin.ProjectDomainAttributesUpdateResponse, error)
}

