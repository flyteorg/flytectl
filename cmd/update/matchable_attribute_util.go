package update

import (
	"context"
	"encoding/json"
	"fmt"

	sconfig "github.com/flyteorg/flytectl/cmd/config/subcommand"
	cmdCore "github.com/flyteorg/flytectl/cmd/core"
	"github.com/flyteorg/flytectl/pkg/ext"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/admin"
)

func DecorateAndUpdateMatchableAttr(
	ctx context.Context,
	cmdCtx cmdCore.CommandContext,
	project, domain, workflow string,
	resourceType admin.MatchableResource,
	attributeDecorator sconfig.MatchableAttributeDecorator,
	dryRun bool,
) error {
	if project == "" {
		return fmt.Errorf("project is required")
	}
	if domain == "" && workflow != "" {
		return fmt.Errorf("domain is required")
	}

	switch {
	case workflow != "":
		return updateWorkflowMatchableAttributes(ctx, cmdCtx, project, domain, workflow, resourceType, attributeDecorator, dryRun)
	case domain != "":
		return updateProjectDomainMatchableAttributes(ctx, cmdCtx, project, domain, resourceType, attributeDecorator, dryRun)
	default:
		return updateProjectMatchableAttributes(ctx, cmdCtx, project, resourceType, attributeDecorator, dryRun)
	}
}

func updateProjectMatchableAttributes(
	ctx context.Context,
	cmdCtx cmdCore.CommandContext,
	project string,
	resourceType admin.MatchableResource,
	attributeDecorator sconfig.MatchableAttributeDecorator,
	dryRun bool,
) error {
	if project == "" {
		panic("project is empty")
	}

	response, err := cmdCtx.AdminFetcherExt().FetchProjectAttributes(ctx, project, resourceType)
	if err != nil && !ext.IsNotFoundError(err) {
		return err
	}

	oldMatchingAttributes := response.GetAttributes().GetMatchingAttributes()
	newMatchingAttributes := attributeDecorator.Decorate()

	v, _ := json.MarshalIndent(oldMatchingAttributes, "", "    ")
	fmt.Println(string(v))

	if dryRun {
		fmt.Printf("Skipping UpdateProjectAttributes request (dryRun)\n")
		return nil
	}

	// TODO: kamal - ack/force

	if err := cmdCtx.AdminUpdaterExt().UpdateProjectAttributes(ctx, project, newMatchingAttributes); err != nil {
		return err
	}

	fmt.Printf("Updated attributes from %s project\n", project)
	return nil
}

func updateProjectDomainMatchableAttributes(
	ctx context.Context,
	cmdCtx cmdCore.CommandContext,
	project, domain string,
	resourceType admin.MatchableResource,
	attributeDecorator sconfig.MatchableAttributeDecorator,
	dryRun bool,
) error {
	if project == "" {
		panic("project is empty")
	}
	if domain == "" {
		panic("domain is empty")
	}

	response, err := cmdCtx.AdminFetcherExt().FetchProjectDomainAttributes(ctx, project, domain, resourceType)
	if err != nil && !ext.IsNotFoundError(err) {
		return err
	}

	oldMatchingAttributes := response.GetAttributes().GetMatchingAttributes()
	newMatchingAttributes := attributeDecorator.Decorate()

	v, _ := json.MarshalIndent(oldMatchingAttributes, "", "    ")
	fmt.Println(string(v))

	if dryRun {
		fmt.Printf("Skipping UpdateProjectDomainAttributes request (dryRun)\n")
		return nil
	}

	// TODO: kamal - ack/force

	if err := cmdCtx.AdminUpdaterExt().UpdateProjectDomainAttributes(ctx, project, domain, newMatchingAttributes); err != nil {
		return err
	}

	fmt.Printf("Updated attributes from %s project and domain %s\n", project, domain)
	return nil
}

func updateWorkflowMatchableAttributes(
	ctx context.Context,
	cmdCtx cmdCore.CommandContext,
	project, domain, workflow string,
	resourceType admin.MatchableResource,
	attributeDecorator sconfig.MatchableAttributeDecorator,
	dryRun bool,
) error {
	if project == "" {
		panic("project is empty")
	}
	if domain == "" {
		panic("domain is empty")
	}
	if workflow == "" {
		panic("workflow is empty")
	}

	response, err := cmdCtx.AdminFetcherExt().FetchWorkflowAttributes(ctx, project, domain, workflow, resourceType)
	if err != nil && !ext.IsNotFoundError(err) {
		return err
	}

	oldMatchingAttributes := response.GetAttributes().GetMatchingAttributes()
	newMatchingAttributes := attributeDecorator.Decorate()

	v, _ := json.MarshalIndent(oldMatchingAttributes, "", "    ")
	fmt.Println(string(v))

	if dryRun {
		fmt.Printf("Skipping UpdateWorkflowAttributes request (dryRun)\n")
		return nil
	}

	// TODO: kamal - ack/force

	if err := cmdCtx.AdminUpdaterExt().UpdateWorkflowAttributes(ctx, project, domain, workflow, newMatchingAttributes); err != nil {
		return err
	}

	fmt.Printf("Updated attributes from %s project and domain %s and workflow %s\n", project, domain, workflow)
	return nil
}
