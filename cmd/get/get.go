package get

import (
	"context"
	"fmt"
	"github.com/lyft/flytectl/cmd/config"
	"github.com/lyft/flytectl/cmd/core"

	"github.com/lyft/flytestdlib/logger"

	"github.com/landoop/tableprinter"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/spf13/cobra"
)

func CreateGetCommand() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve various resource.",
	}

	getResourcesFuncs := map[string]core.CommandFunc{
		"projects": getProjectsFunc,
		"domains":  getDomainFuc,
		"tasks":    getTaskFunc,
		"workflows":    getWorkflowFunc,
	}

	core.AddCommands(getCmd, getResourcesFuncs)

	return getCmd
}

func getProjectsFunc(ctx context.Context, args []string, cmdCtx core.CommandContext) error {
	projects, err := cmdCtx.AdminClient().ListProjects(ctx, &admin.ProjectListRequest{})
	if err != nil {
		return err
	}
	logger.Debugf(ctx, "Retrieved %v projects", len(projects.Projects))
	printer := tableprinter.New(cmdCtx.OutputPipe())
	printer.Print(toPrintableProjects(projects.Projects))
	return nil
}

func getTaskFunc(ctx context.Context, args []string, cmdCtx core.CommandContext) error {
	tasks, err := cmdCtx.AdminClient().ListTaskIds(ctx, &admin.NamedEntityIdentifierListRequest{
		Project: config.GetConfig().Project,
		Domain:  config.GetConfig().Domain,
		Limit:   3,
	})
	if err != nil {
		return err
	}

	logger.Debugf(ctx, "Retrieved %v Task", len(tasks.Entities))
	printer := tableprinter.New(cmdCtx.OutputPipe())
	printer.Print(toPrintableTask(tasks.Entities))
	return nil
}

func getWorkflowFunc(ctx context.Context, args []string, cmdCtx core.CommandContext) error {
	if config.GetConfig().Project == "" {
		return fmt.Errorf("Please set project name to get domain")
	}
	if config.GetConfig().Domain == "" {
		return fmt.Errorf("Please set project name to get workflow")
	}
	if len(args) > 0 {
		//workflows, err := cmdCtx.AdminClient().GetWorkflow(ctx, &admin.ObjectGetRequest{
		//	Id : args[0],
		//})
	}
	workflows, err := cmdCtx.AdminClient().ListWorkflowIds(ctx, &admin.NamedEntityIdentifierListRequest{
		Project: config.GetConfig().Project,
		Domain:  config.GetConfig().Domain,
		Limit: 3,
	})
	if err != nil {
		return err
	}
	logger.Debugf(ctx, "Retrieved %v workflows", len(workflows.Entities))
	printer := tableprinter.New(cmdCtx.OutputPipe())
	printer.Print(toPrintableWorkflow(workflows.Entities))
	return nil
}

func getDomainFuc(ctx context.Context, args []string, cmdCtx core.CommandContext) error {
	if config.GetConfig().Project == "" {
		return fmt.Errorf("Please set project name to get domain")
	}
	projects, err := cmdCtx.AdminClient().ListProjects(ctx, &admin.ProjectListRequest{})
	if err != nil {
		return err
	}
	logger.Debugf(ctx, "Retrieved %v domain", len(projects.Projects))
	printer := tableprinter.New(cmdCtx.OutputPipe())
	printer.Print(toPrintableDomain(projects.Projects))
	return nil
}

func toPrintableProjects(projects []*admin.Project) []interface{} {
	type printableProject struct {
		Id          string `header:"Id"`
		Name        string `header:"Name"`
		Description string `header:"Description"`
	}

	res := make([]interface{}, 0, len(projects))
	for _, p := range projects {
		res = append(res, printableProject{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
		})
	}

	return res
}

func toPrintableTask(tasks []*admin.NamedEntityIdentifier) []interface{} {
	type printableTask struct {
		Name    string `header:"Name"`
		Project string `header:"Project"`
		Domain  string `header:"Domain"`
	}

	res := make([]interface{}, 0, len(tasks))
	for _, p := range tasks {
		res = append(res, printableTask{
			Name:    p.Name,
			Domain:  p.Domain,
			Project: p.Project,
		})
	}

	return res
}

func toPrintableDomain(projects []*admin.Project) []interface{} {
	type printableDomain struct {
		Name string `header:"Name"`
		Id   string `header:"Id"`
	}

	res := make([]interface{}, 0, len(projects))
	for _, p := range projects {
		if p.Name == config.GetConfig().Project {
			for _, d := range p.Domains {
				res = append(res, printableDomain{
					Id:   d.Id,
					Name: d.Name,
				})
			}
			break
		}
	}

	return res
}

func toPrintableWorkflow(workflows []*admin.NamedEntityIdentifier) []interface{} {
	type printableWorkflow struct {
		Name string `header:"Name"`
		Domain   string `header:"Domain"`
		Project   string `header:"Project"`
	}

	res := make([]interface{}, 0, len(workflows))
	for _, w := range workflows {
		res = append(res, printableWorkflow{
			Name:   w.Name,
			Domain: w.Domain,
			Project: w.Project,
		})
	}

	return res
}