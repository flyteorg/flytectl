package printer

import (
	"encoding/json"
	"fmt"
	"github.com/landoop/tableprinter"
	"github.com/yalp/jsonpath"
	"os"
)

type Printer struct{}

func (p Printer) PrintOutput(output string, i interface{}) {
	// Factory Method for all printer
	fmt.Println("==",output)
	switch output {
	case "json": // Print protobuf to json
		result, err := json.Marshal(i)
		if err != nil {
			os.Exit(1)
		}
		fmt.Println(string(result))
		break
	case "yaml": // Print protobuf to yaml
	    
		break
	default: // Print table

		printer := tableprinter.New(os.Stdout)
		printer.Print(i)
		break
	}
}

func(p Printer) BuildOutput(input []interface{},column map[string]string,printTransform func(data []byte)(interface{},error)) ([]interface{},error) {
	responses := make([]interface{}, 0, len(input))
	for _, data := range input {
		tableData := make(map[string]interface{})
		for k := range column {
			data, _ := jsonpath.Read(data, column[k])
			tableData[k] = data.(string)
		}
		jsonbody, err := json.Marshal(tableData)
		if err != nil {
			return responses,err
		}
		response,err := printTransform(jsonbody)
		if err != nil {
			return responses,err
		}
		responses = append(responses, response)
	}
	return responses,nil
}

func (p Printer) Print(output string, i interface{},column map[string]string,printTransform func(data []byte)(interface{},error)) {
	var data interface{}
	byte, _ := json.Marshal(i)
	_ = json.Unmarshal(byte, &data)
	input := data.([]interface{})
	response,err := p.BuildOutput(input,column,printTransform)
	if err != nil {
		os.Exit(1)
	}
	p.PrintOutput(output, response)
}
