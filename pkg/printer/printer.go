package printer

import (
	"encoding/json"
	"github.com/landoop/tableprinter"
	"github.com/yalp/jsonpath"
	"os"
)

type Printer struct{}

func (p Printer) Print(output string, i interface{}) {
	// Factory Method for all printer
	switch output {
	case "json": // Print protobuf to json
		break
	case "yaml": // Print protobuf to yaml
		break
	default: // Print table
		printer := tableprinter.New(os.Stdout)
		printer.Print(i)
		break
	}
}

func (p Printer) PrintBuildNamedEntityIdentifier(output string, i interface{},column map[string]string) error {
	var entity interface{}
	byte, _ := json.Marshal(i)
	_ = json.Unmarshal(byte, &entity)
	obj := entity.([]interface{})
	res := make([]interface{}, 0, len(obj))
	for _, data := range obj {
		var results PrintableNamedEntityIdentifier
		entityTable := make(map[string]interface{})
		for k := range column {
			data, _ := jsonpath.Read(data, column[k])
			entityTable[k] = data.(string)
		}
		jsonbody, err := json.Marshal(entityTable)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(jsonbody, &results); err != nil {
			return err
		}
		res = append(res, results)
	}
	return nil
}

func (p Printer) PrintTask(output string, i interface{},column map[string]string) {
	var task interface{}
	byte, _ := json.Marshal(i)
	_ = json.Unmarshal(byte, &task)
	obj := task.([]interface{})
	res := make([]interface{}, 0, len(obj))
	for _, data := range obj {
			var results PrintableTask
			taskTable := make(map[string]interface{})
			for k := range column {
				data, _ := jsonpath.Read(data, column[k])
				taskTable[k] = data.(string)
			}
			jsonbody, err := json.Marshal(taskTable)
			if err != nil {
				return
			}
			if err := json.Unmarshal(jsonbody, &results); err != nil {
				return
			}
			res = append(res, results)
	}
	p.Print(output, res)
}

func (p Printer) PrintWorkflow(output string, i interface{},column map[string]string) {
	var task []interface{}
	byte, _ := json.Marshal(i)
	_ = json.Unmarshal(byte, &task)
	res := make([]interface{}, 0, len(task))
	for _, data := range task {
			var results PrintableWorkflow
			workflowTable := make(map[string]interface{})
			for k := range column {
				data, _ := jsonpath.Read(data, column[k])
				workflowTable[k] = data.(string)
			}
			jsonbody, err := json.Marshal(workflowTable)
			if err != nil {
				return
			}
			if err := json.Unmarshal(jsonbody, &results); err != nil {
				return
			}
			res = append(res, results)
	}
	p.Print(output, res)
}

func (p Printer) PrintProject(output string, i interface{},column map[string]string) {
	var projects interface{}
	byte, _ := json.Marshal(i)
	_ = json.Unmarshal(byte, &projects)
	obj := projects.([]interface{})
	res := make([]interface{}, 0, len(obj))
	for _, p := range obj {
		var results PrintableProject
		projectTable := make(map[string]interface{})
		for k := range column {
			data, _ := jsonpath.Read(p, column[k])
			projectTable[k] = data.(string)
        }
		jsonbody, err := json.Marshal(projectTable)
		if err != nil {
			return
		}
		if err := json.Unmarshal(jsonbody, &results); err != nil {
			return
		}

		res = append(res, results)
	}
	p.Print(output, res)
}

func (p Printer) PrintDomain(output string, i interface{},column map[string]string) {
	var domains interface{}
	byte, _ := json.Marshal(i)
	_ = json.Unmarshal(byte, &domains)
	obj := domains.([]interface{})
	res := make([]interface{}, 0, len(obj))
	for _, p := range obj {
		var results PrintableDomain
		projectTable := make(map[string]interface{})
		for k := range column {
			data, _ := jsonpath.Read(p, column[k])
			projectTable[k] = data.(string)
		}

		jsonbody, err := json.Marshal(projectTable)
		if err != nil {
			return
		}
		if err := json.Unmarshal(jsonbody, &results); err != nil {
			return
		}

		res = append(res, results)
	}
	p.Print(output, res)
}
