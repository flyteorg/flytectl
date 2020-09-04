package printer

import (
	"encoding/json"
	"github.com/landoop/tableprinter"
	cmdCore "github.com/lyft/flytectl/cmd/core"
	"github.com/lyft/flyteidl/gen/pb-go/flyteidl/admin"
	"github.com/yalp/jsonpath"
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
	var projects interface{}
	byte, _ := json.Marshal(i)
	_ = json.Unmarshal(byte, &projects)
	obj := projects.([]interface{})
	res := make([]interface{}, 0, len(obj))
	for _, p := range obj {
		id, _ := jsonpath.Read(p, "$.id")
		name, _ := jsonpath.Read(p, "$.name")
		description, _ := jsonpath.Read(p, "$.description")
			res = append(res, PrintableProject{
				Id:          id.(string),
				Name:        name.(string),
				Description: description.(string),
			})
	}

	p.Printer.Print(output, res, p.Ctx)
}

type DomainList struct {
	Ctx cmdCore.CommandContext
	Printer
}

func (d DomainList) Print(output string, i interface{}) {
	var domains interface{}
	byte, _ := json.Marshal(i)
	_ = json.Unmarshal(byte, &domains)
	obj := domains.([]interface{})
	res := make([]interface{}, 0, len(obj))
	for _, domain := range obj {
		id, _ := jsonpath.Read(domain, "$.id")
		name, _ := jsonpath.Read(domain, "$.name")
			res = append(res, PrintableDomain{
				Id:   id.(string),
				Name: name.(string),
			})
	}
	d.Printer.Print(output, res, d.Ctx)
}
