package printer

import (
	"github.com/landoop/tableprinter"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
)

type Printer struct{}

func (p Printer) Print(output string, i interface{}, cmdCtx cmdCore.CommandContext) {
	// Factory Method for all printer
	switch output {
	case "json": // Print protobuf to json
		break
	case "yaml": // Print protobuf to yaml
		break
	default: // Print table
		printer := tableprinter.New(cmdCtx.OutputPipe())
		printer.Print(i)
		break
	}
}

func (p Printer) buildNamedEntityIdentifier(data *admin.NamedEntityIdentifier) interface{} {
	return PrintableNamedEntityIdentifier{
		Name:    data.Name,
		Domain:  data.Domain,
		Project: data.Project,
	}
}

type AdminTasksList struct {
	Ctx cmdCore.CommandContext
	Printer
}

func (a AdminTasksList) Print(output string, i interface{}) {
	input := i.(map[int]interface{})
	res := make([]interface{}, 0, len(input))
	for _, data := range input {
		switch data.(type) {
		case *admin.Task:
			task := data.(*admin.Task)
			res = append(res, PrintableTask{
				Version:          task.Id.Version,
				Name:             task.Id.Name,
				Type:             task.Closure.CompiledTask.Template.Type,
				Discoverable:     task.Closure.CompiledTask.Template.Metadata.Discoverable,
				DiscoveryVersion: task.Closure.CompiledTask.Template.Metadata.DiscoveryVersion,
			})
			break
		case *admin.NamedEntityIdentifier:
			res = append(res, a.buildNamedEntityIdentifier(i.(*admin.NamedEntityIdentifier)))

			break
		}
	}
	a.Printer.Print(output, res, a.Ctx)
}

type AdminWorkflowsList struct {
	Ctx cmdCore.CommandContext
	Printer
}

func (a AdminWorkflowsList) Print(output string, i interface{}) {
	input := i.(map[int]interface{})
	res := make([]interface{}, 0, len(input))
	for _, data := range input {
		switch data.(type) {
		case *admin.Workflow:
			workflow := data.(*admin.Workflow)
			res = append(res, PrintableWorkflow{
				Version: workflow.Id.Version,
				Name:    workflow.Id.Name,
			})
			break
		case *admin.NamedEntityIdentifier:
			res = append(res, a.buildNamedEntityIdentifier(i.(*admin.NamedEntityIdentifier)))
			break
		}
	}
	a.Printer.Print(output, res, a.Ctx)
}

type ProjectList struct {
	Ctx cmdCore.CommandContext
	Printer
}

func (p ProjectList) Print(output string, i interface{}) {
	projects := i.(map[int]interface{})
	res := make([]interface{}, 0, len(projects))

	for _, project := range projects {
		switch project.(type) {
		case *admin.Project:
			p := project.(*admin.Project)
			res = append(res, PrintableProject{
				Id:          p.Id,
				Name:        p.Name,
				Description: p.Description,
			})
			break
		}
	}
	p.Printer.Print(output, res, p.Ctx)
}

type DomainList struct {
	Ctx cmdCore.CommandContext
	Printer
}

func (d DomainList) Print(output string, i interface{}) {
	domains := i.(map[int]interface{})
	res := make([]interface{}, 0, len(domains))

	for _, domain := range domains {
		switch domain.(type) {
		case *admin.Project:
			p := domain.(*admin.Domain)
			res = append(res, PrintableDomain{
				Id:   p.Id,
				Name: p.Name,
			})
			break
		}
	}
	d.Printer.Print(output, res, d.Ctx)
}
