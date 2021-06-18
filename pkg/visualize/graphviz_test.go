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
			fmt.Println(b)
			assert.NoError(t, err)
			assert.NotNil(t, b)
		})
	}
}

func TestAddBranchSubNodeEdge(t *testing.T) {
	//gb := newGraphBuilder()
	//gb.nodeClusters["n"] = "innerGraph"
	//mockGraph := &mocks.Graphvizer{}
	//attrs := map[string]string{}
	//attrs[LHeadAttr] = "innerGraph"
	//attrs[LabelAttr] = fmt.Sprintf("\"%s\"", "label")
	//// Verify the attributes
	//mockGraph.OnAddEdgeMatch(mock.Anything, mock.Anything, mock.Anything, attrs).Return(nil)
	//mockGraph.OnGetEdgeMatch(mock.Anything, mock.Anything).Return(&graphviz.Edge{})
	//parentNode := &graphviz.Node{Name : "parentNode", Attrs: nil}
	//n := &graphviz.Node{Name: "n"}
	////err := gb.addBranchSubNodeEdge(mockGraph, parentNode, n, "label")
	//assert.NoError(t, err)
}
