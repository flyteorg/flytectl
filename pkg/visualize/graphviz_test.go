package visualize

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/golang/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
)

func TestRenderWorkflowBranch(t *testing.T) {
	// Sadly we cannot compare the output of svg, as it slightly changes.
	file := []string{"compiled_closure_branch_nested", "compiled_subworkflows"}

	for _, s := range file {
		t.Run(s, func(t *testing.T) {
			r, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", s))
			assert.NoError(t, err)

			i := bytes.NewReader(r)

			c := &core.CompiledWorkflowClosure{}
			err = jsonpb.Unmarshal(i, c)
			assert.NoError(t, err)
			b, err := RenderWorkflow(c)
			assert.NoError(t, err)
			assert.NotNil(t, b)
		})
	}
}
