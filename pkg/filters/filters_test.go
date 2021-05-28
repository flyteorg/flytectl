package filters

import (
	"fmt"
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
			Input:  "project.name=flytesnacks",
			Output: "eq(project.name,flytesnacks)",
		},
		{
			Input:  "execution.phase in (FAILED;SUCCEEDED),execution.name=y8n2wtuspj",
			Output: "value_in(execution.phase,FAILED;SUCCEEDED)+eq(execution.name,y8n2wtuspj)",
		},
	}
	for _, test := range tests {
		filters := SplitTerms(test.Input)

		result, err := Transform(filters)
		fmt.Println(result)
		assert.Nil(t, err)
		assert.Equal(t, test.Output, result)
	}
}
