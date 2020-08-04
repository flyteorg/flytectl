package get

import (
	"context"
	"fmt"
	"github.com/lyft/flytectl/cmd/config"
	"github.com/lyft/flytectl/cmd/core"

	core "github.com/lyft/flyteidl/gen/pb-go/flyteidl/core"
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

	getResourcesFuncs := map[string]cmdCore.CommandFunc{
		"projects": getProjectsFunc,
		"domains":  getDomainFunc,
		"tasks":    getTaskFunc,
		"workflows":    getWorkflowFunc,
	}

	cmdCore.AddCommands(getCmd, getResourcesFuncs)

	return getCmd
}

func getProjectsFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	projects, err := cmdCtx.AdminClient().ListProjects(ctx, &admin.ProjectListRequest{})
	if err != nil {
		return err
	}
	logger.Debugf(ctx, "Retrieved %v projects", len(projects.Projects))
	printer := tableprinter.New(cmdCtx.OutputPipe())
	printer.Print(toPrintableProjects(projects.Projects))
	return nil
}

func getTaskFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	if config.GetConfig().Project == "" {
		return fmt.Errorf("Please set project name to get domain")
	}
	if config.GetConfig().Domain == "" {
		return fmt.Errorf("Please set project name to get workflow")
	}
	if len(args) == 1 {
		task, err := cmdCtx.AdminClient().ListTasks(ctx, &admin.ResourceListRequest{
			Id : &admin.NamedEntityIdentifier{
				Project: config.GetConfig().Project,
				Domain:  config.GetConfig().Domain,
				Name: args[0],
			},
			Limit:   3,
		})
		if err != nil {
			return err
		}
		logger.Debugf(ctx, "Retrieved Task",)
		printer := tableprinter.New(cmdCtx.OutputPipe())
		printer.Print(toPrintableGetTask(task.Tasks))
		return nil
	}

	tasks, err := cmdCtx.AdminClient().ListTaskIds(ctx, &admin.NamedEntityIdentifierListRequest{
		Project: config.GetConfig().Project,
		Domain:  config.GetConfig().Domain,
		Limit: 10,
	})
	if err != nil {
		return err
	}
	logger.Debugf(ctx, "Retrieved %v Task", len(tasks.Entities))

	printer := tableprinter.New(cmdCtx.OutputPipe())
	printer.Print(toPrintableTask(tasks.Entities))
	return nil
}

func getWorkflowFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
	if config.GetConfig().Project == "" {
		return fmt.Errorf("Please set project name to get domain")
	}
	if config.GetConfig().Domain == "" {
		return fmt.Errorf("Please set project name to get workflow")
	}
	if len(args) > 0 {
		workflows, err := cmdCtx.AdminClient().ListWorkflows(ctx, &admin.ResourceListRequest{
			Id : &admin.NamedEntityIdentifier{
				Project: config.GetConfig().Project,
				Domain:  config.GetConfig().Domain,
				Name: args[0],
			},
			Limit:   3,
		})
		if err != nil {
			return err
		}
		logger.Debugf(ctx, "Retrieved %v workflows", len(workflows.Workflows))
		printer := tableprinter.New(cmdCtx.OutputPipe())
		printer.Print(toPrintableListWorkflow(workflows.Workflows))
		return nil
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

func getDomainFunc(ctx context.Context, args []string, cmdCtx cmdCore.CommandContext) error {
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
	type printableTasks struct {
		Name    string `header:"Name"`
		Project string `header:"Project"`
		Domain  string `header:"Domain"`
	}

	res := make([]interface{}, 0, len(tasks))
	for _, p := range tasks {
		res = append(res, printableTasks{
			Name:    p.Name,
			Domain:  p.Domain,
			Project: p.Project,
		})
	}

	return res
}

func toPrintableGetTask(tasks []*admin.Task) []interface{} {
	type printableTask struct {
		Version    string `header:"Version"`
		Name    string `header:"Name"`
		Type  string `header:"Type"`
		Discoverable  bool `header:"Discoverable"`
		ResourceType core.ResourceType `header:"ResourceType"`
	}
	res := make([]interface{}, 0, len(tasks))
	for _,task := range tasks {
		res = append(res,printableTask{
		Version:    task.Id.Version,
			Name: task.Id.Name,
				ResourceType: task.Id.ResourceType,
				Type : task.Closure.CompiledTask.Template.Type,
				Discoverable : task.Closure.CompiledTask.Template.Metadata.Discoverable,
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

func toPrintableListWorkflow(workflows []*admin.Workflow) []interface{} {
	type printableWorkflow struct {
		Version    string `header:"Version"`
		Name   string `header:"Name"`
		ResourceType core.ResourceType `header:"ResourceType"`
	}
	res := make([]interface{}, 0, len(workflows))
	for _,workflow := range workflows {
		res = append(res,printableWorkflow{
			Version:    workflow.Id.Version,
			Name: workflow.Id.Name,
			ResourceType: workflow.Id.ResourceType,
		})
	}
	return res
}