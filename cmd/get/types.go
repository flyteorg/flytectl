package get

import "encoding/json"

// PrintableWorkflow is the structure for printing workflow
type PrintableWorkflow struct {
	Name    string `header:"Name"`
	Version string `header:"Version"`
}

// PrintableTask is the structure for printing Task
type PrintableTask struct {
	Version          string `header:"Version"`
	Name             string `header:"Name"`
	Type             string `header:"Type"`
	Discoverable     bool   `header:"Discoverable"`
	DiscoveryVersion string `header:"DiscoveryVersion"`
}

// PrintableProject is the structure for printing Project
type PrintableProject struct {
	ID          string `header:"Id"`
	Name        string `header:"Name"`
	Description string `header:"Description"`
}

// PrintableSingleExecution is the structure for printing Execution
type PrintableSingleExecution struct {
	Version    string `header:"Version"`
	Name       string `header:"Name"`
	Type       string `header:"Type"`
	LaunchPlan string `header:"LaunchPlan"`
	Phase      string `header:"Phase"`
	Duration   string `header:"Duration"`
	Workflow   string `header:"Workflow"`
	Metadata   string `header:"Metadata"`
}

// PrintableLaunchPlan is the structure for printing Launch Plan
type PrintableLaunchPlan struct {
	Version string `header:"Version"`
	Name    string `header:"Name"`
	Type    string `header:"Type"`
}

// PrintableExecution is the structure for printing Excution
type PrintableExecution struct {
	Version string `header:"Version"`
	Name    string `header:"Name"`
	Type    string `header:"Type"`
}

// PrintableNamedEntityIdentifier
type PrintableNamedEntityIdentifier struct {
	Name    string `header:"Name"`
	Project string `header:"Project"`
	Domain  string `header:"Domain"`
}

var entityStructure = map[string]string{
	"Domain":  "$.domain",
	"Name":    "$.name",
	"Project": "$.project",
}

var transformTaskEntity = func(jsonbody []byte) (interface{}, error) {
	results := PrintableNamedEntityIdentifier{}
	if err := json.Unmarshal(jsonbody, &results); err != nil {
		return results, err
	}
	return results, nil
}
