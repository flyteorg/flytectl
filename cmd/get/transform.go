package get

import (
	"encoding/json"
)

func transformExecution(jsonbody []byte) (interface{}, error) {
	results := PrintableExecution{}
	if err := json.Unmarshal(jsonbody, &results); err != nil {
		return results, err
	}
	return results, nil
}

func transformSingleExecution(jsonbody []byte) (interface{}, error) {
	results := PrintableSingleExecution{}
	if err := json.Unmarshal(jsonbody, &results); err != nil {
		return results, err
	}
	return results, nil
}

func transformLaunchPlan(jsonbody []byte) (interface{}, error) {
	results := PrintableLaunchPlan{}
	if err := json.Unmarshal(jsonbody, &results); err != nil {
		return results, err
	}
	return results, nil
}

func transformProject(jsonbody []byte) (interface{}, error) {
	results := PrintableProject{}
	if err := json.Unmarshal(jsonbody, &results); err != nil {
		return results, err
	}
	return results, nil
}

var transformTask = func(jsonbody []byte) (interface{}, error) {
	results := PrintableTask{}
	if err := json.Unmarshal(jsonbody, &results); err != nil {
		return results, err
	}
	return results, nil
}

var transformWorkflow = func(jsonbody []byte) (interface{}, error) {
	results := PrintableWorkflow{}
	if err := json.Unmarshal(jsonbody, &results); err != nil {
		return results, err
	}
	return results, nil
}
