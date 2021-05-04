package interfaces

import (
	"context"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

//go:generate mockery -all -case=underscore

// Fetcher Interface for exposing the fetch capabilities from the admin and also allow this to be injectable into other
// modules. eg : create execution which requires to fetch launchplan details.
type Fetcher interface {
	FetchExecution(ctx context.Context, adminClient service.AdminServiceClient, name, project,
		domain string) (*admin.Execution, error)

	FetchLPForName(ctx context.Context, adminClient service.AdminServiceClient, name, project,
		domain string) ([]*admin.LaunchPlan, error)

	FetchAllVerOfLP(ctx context.Context, adminClient service.AdminServiceClient, lpName, project,
		domain string) ([]*admin.LaunchPlan, error)

	FetchLPLatestVersion(ctx context.Context, adminClient service.AdminServiceClient, name, project,
		domain string) (*admin.LaunchPlan, error)

	FetchLPVersion(ctx context.Context, adminClient service.AdminServiceClient, name, version, project,
		domain string)(*admin.LaunchPlan, error)

	GetWorkflowAttributes(ctx context.Context, adminClient service.AdminServiceClient, project, domain, name string,
		rsType admin.MatchableResource) (*admin.WorkflowAttributesGetResponse, error)

	GetProjectDomainAttributes(ctx context.Context, adminClient service.AdminServiceClient, project, domain string,
		rsType admin.MatchableResource) (*admin.ProjectDomainAttributesGetResponse, error)
}
