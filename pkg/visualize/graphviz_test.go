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

func TestRenderWorkflow(t *testing.T) {
	r, err := ioutil.ReadFile("testdata/compiled_closure_branch_nested.json")
	assert.NoError(t, err)

	i := bytes.NewReader(r)

	c := &core.CompiledWorkflowClosure{}
	err = jsonpb.Unmarshal(i, c)
	assert.NoError(t, err)
	b, err := RenderWorkflow(c)
	assert.NoError(t, err)
	assert.NotNil(t, b)
	fmt.Println(string(b))
}
