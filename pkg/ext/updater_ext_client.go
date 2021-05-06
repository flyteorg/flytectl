package ext

import (
	"context"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/service"
)

//go:generate mockery -all -case=underscore

// AdminUpdaterExtInterface Interface for exposing the update capabilities from the admin
type AdminUpdaterExtInterface interface {
	AdminServiceClient() service.AdminServiceClient

	// UpdateWorkflowAttributes fetches workflow attributes within a project, domain for a particular matchable resource
	UpdateWorkflowAttributes(ctx context.Context, project, domain, name string, rsType admin.MatchableResource) error

	// UpdateProjectDomainAttributes fetches project domain attributes for a particular matchable resource
	UpdateProjectDomainAttributes(ctx context.Context, project, domain string, rsType admin.MatchableResource) error
}

// AdminUpdaterExtClient is used for interacting with extended features used for updating data in admin service
type AdminUpdaterExtClient struct {
	AdminClient service.AdminServiceClient
}

func (a *AdminUpdaterExtClient) UpdateWorkflowAttributes(ctx context.Context, project, domain, name string, rsType admin.MatchableResource) error {
	_, err := a.AdminServiceClient().UpdateWorkflowAttributes(ctx, &admin.WorkflowAttributesUpdateRequest{
		Attributes: &admin.WorkflowAttributes{
			Project: project,
			Domain: domain,
			Workflow: name,
			MatchingAttributes: &admin.MatchingAttributes{
				Target: &admin.MatchingAttributes_TaskResourceAttributes{
					TaskResourceAttributes: &admin.TaskResourceAttributes{
						Defaults: &admin.TaskResourceSpec{
							Cpu: "1",
							Memory: "150Mi",
						},
					},
				},
			},
		},
	})
	return err
}

func (a *AdminUpdaterExtClient) UpdateProjectDomainAttributes(ctx context.Context, project, domain string, rsType admin.MatchableResource) error {
	_, err := a.AdminServiceClient().UpdateProjectDomainAttributes(ctx,
		&admin.ProjectDomainAttributesUpdateRequest{
			Attributes:&admin.ProjectDomainAttributes{
				Project: project,
				Domain: domain,
				MatchingAttributes:  &admin.MatchingAttributes{
					Target: &admin.MatchingAttributes_TaskResourceAttributes{
						TaskResourceAttributes: &admin.TaskResourceAttributes{
							Defaults: &admin.TaskResourceSpec{
								Cpu: "2",
								Memory: "200Mi",
							},
						},
					},
				},
			},
		})
	return err
}

func (a *AdminUpdaterExtClient) AdminServiceClient() service.AdminServiceClient {
	if a == nil {
		return nil
	}
	return a.AdminClient
}
