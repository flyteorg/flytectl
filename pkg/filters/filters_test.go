package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

func TestTransformFilter(t *testing.T) {
	tests := []TestCase{
		{
			Input:  "project.Value>4,project.Value<4",
			Output: "gt(project.Value,4)+lt(project.Value,4)",
		},
		{
			Input:  "project.Phase in (RUNNING;SUCCESS),project.Phase contains RUNNING",
			Output: "value_in(project.Phase,RUNNING;SUCCESS)+contains(project.Phase,RUNNING)",
		},
	}
	for _, test := range tests {
		filters := SplitTerms(test.Input)

		result, err := Transform(filters)
		assert.Nil(t, err)
		assert.Equal(t, test.Output, result)
	}
}
