package ext

import (
	"context"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
)

//go:generate mockery -all -case=underscore

// AdminDeleterExtInterface Interface for exposing the update capabilities from the admin
type AdminDeleterExtInterface interface {
	AdminServiceClient() service.AdminServiceClient

	// DeleteWorkflowAttributes fetches workflow attributes within a project, domain for a particular matchable resource
	DeleteWorkflowAttributes(ctx context.Context, project, domain, name string, rsType admin.MatchableResource) error

	// UpdateProjectDomainAttributes fetches project domain attributes for a particular matchable resource
	DeleteProjectDomainAttributes(ctx context.Context, project, domain string, rsType admin.MatchableResource) error
}

// AdminDeleterExtClient is used for interacting with extended features used for deleting/archiving data in admin service
type AdminDeleterExtClient struct {
	AdminClient service.AdminServiceClient
}

func (a *AdminDeleterExtClient) DeleteWorkflowAttributes(ctx context.Context, project, domain, name string, rsType admin.MatchableResource) error {
	_, err := a.AdminServiceClient().DeleteWorkflowAttributes(ctx, &admin.WorkflowAttributesDeleteRequest{
		Project: project,
		Domain: domain,
		Workflow: name,
		ResourceType: rsType,
	})
	return err
}

func (a *AdminDeleterExtClient) DeleteProjectDomainAttributes(ctx context.Context, project, domain string, rsType admin.MatchableResource) error {
	_, err := a.AdminServiceClient().DeleteProjectDomainAttributes(ctx, &admin.ProjectDomainAttributesDeleteRequest{
		Project: project,
		Domain: domain,
		ResourceType: rsType,
	})
	return err
}

func (a *AdminDeleterExtClient) AdminServiceClient() service.AdminServiceClient {
	if a == nil {
		return nil
	}
	return a.AdminClient
}
