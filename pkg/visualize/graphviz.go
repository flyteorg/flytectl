package visualize

import (
	"bytes"
	"context"
	"fmt"
	"github.com/flyteorg/flyteidl/clients/go/coreutils"
	"strings"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/flyteorg/flytestdlib/errors"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func operandToString(op *core.Operand) string {
	if op.GetPrimitive() != nil {
		l, err := coreutils.ExtractFromLiteral(&core.Literal{Value: &core.Literal_Scalar{
			Scalar: &core.Scalar{
				Value: &core.Scalar_Primitive{
					Primitive: op.GetPrimitive(),
				},
			},
		}})
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("%v", l)
	}
	return op.GetVar()
}

func comparisonToString(expr *core.ComparisonExpression) string {
	return fmt.Sprintf("%s %s %s", operandToString(expr.LeftValue), expr.Operator.String(), operandToString(expr.RightValue))
}

func conjunctionToString(expr *core.ConjunctionExpression) string {
	return fmt.Sprintf("(%s) %s (%s)", booleanExprToString(expr.LeftExpression), expr.Operator.String(), booleanExprToString(expr.RightExpression))
}

func booleanExprToString(expr *core.BooleanExpression) string {
	if expr.GetConjunction() != nil {
		return conjunctionToString(expr.GetConjunction())
	}
	return comparisonToString(expr.GetComparison())
}

func constructStartNode(graph *cgraph.Graph) (*cgraph.Node, error) {
	gn, err := graph.CreateNode("start-node")
	if err != nil {
		return nil, err
	}
	gn.SetLabel("start")
	gn.SetShape(cgraph.DoubleCircleShape)
	gn.SetColor("green")
	return gn, nil
}

func constructEndNode(graph *cgraph.Graph) (*cgraph.Node, error) {
	gn, err := graph.CreateNode("end-node")
	if err != nil {
		return nil, err
	}
	gn.SetLabel("end")
	gn.SetShape(cgraph.DoubleCircleShape)
	gn.SetColor("red")
	return gn, nil
}

func constructTaskNode(name string, graph *cgraph.Graph, n *core.Node, t *core.CompiledTask) (*cgraph.Node, error) {
	gn, err := graph.CreateNode(name)
	if err != nil {
		return nil, err
	}
	if n.Metadata != nil && n.Metadata.Name != "" {
		v := strings.LastIndexAny(n.Metadata.Name, ".")
		gn.SetLabel(fmt.Sprintf("%s [%s]", n.Metadata.Name[v+1:], t.Template.Type))
	}
	gn.SetShape(cgraph.BoxShape)
	return gn, nil
}

func constructErrorNode(name string, graph *cgraph.Graph, m string) (*cgraph.Node, error) {
	gn, err := graph.CreateNode(name)
	if err != nil {
		return nil, err
	}
	gn.SetLabel(m)
	gn.SetShape(cgraph.BoxShape)
	gn.SetColor("red")
	return gn, nil
}

func constructBranchConditionNode(name string, graph *cgraph.Graph, n *core.Node) (*cgraph.Node, error) {
	gn, err := graph.CreateNode(name)
	if err != nil {
		return nil, err
	}
	if n.Metadata != nil && n.Metadata.Name != "" {
		gn.SetLabel(n.Metadata.Name)
	}
	gn.SetShape(cgraph.DiamondShape)
	return gn, nil
}

func getName(prefix, id string) string {
	if prefix != "" {
		return prefix + "-" + id
	}
	return id
}

type graphBuilder struct {
	// Mutated as graph is built
	graphNodes map[string]*cgraph.Node
	// Mutated as graph is built. lookup table for all graphviz compiled edges.
	graphEdges map[string]*cgraph.Edge
	// lookup table for all graphviz compiled subgraphs
	subWf      map[string]*cgraph.Graph
	// a lookup table for all tasks in the graph
	tasks      map[string]*core.CompiledTask
	// a lookup for all node clusters. This is to remap the edges to the cluster itself (instead of the node)
	// this is useful in the case of branchNodes and subworkflow nodes
	nodeClusters map[string]*cgraph.Graph
}

func (gb *graphBuilder) addBranchSubNodeEdge(graph *cgraph.Graph, parentNode, n *cgraph.Node, label string) error {
	edgeName := fmt.Sprintf("%s-%s", parentNode.Name(), n.Name())
	if _, ok := gb.graphEdges[edgeName]; !ok {
		edge, err := graph.CreateEdge(edgeName, parentNode, n)
		if err != nil {
			return err
		}
		edge.SetLabel(label)
		if c, ok := gb.nodeClusters[n.Name()]; ok {
			edge.SetLogicalHead(c.Name())
		}
		gb.graphEdges[edgeName] = edge
	}
	return nil
}

func (gb *graphBuilder) constructBranchNode(prefix string, graph *cgraph.Graph, n *core.Node) (*cgraph.Node, error) {
	parentBranchNode, err := constructBranchConditionNode(getName(prefix, n.Id), graph, n)
	if err != nil {
		return nil, err
	}
	gb.graphNodes[parentBranchNode.Name()] = parentBranchNode

	if n.GetBranchNode().GetIfElse() == nil {
		return parentBranchNode, nil
	}

	subNode, err := gb.constructNode(prefix, graph, n.GetBranchNode().GetIfElse().Case.ThenNode)
	if err != nil {
		return nil, err
	}
	if err := gb.addBranchSubNodeEdge(graph, parentBranchNode, subNode, booleanExprToString(n.GetBranchNode().GetIfElse().Case.Condition)); err != nil {
		return nil, err
	}

	if n.GetBranchNode().GetIfElse().GetError() != nil {
		name := fmt.Sprintf("%s-error", parentBranchNode.Name())
		subNode, err := constructErrorNode(name, graph, n.GetBranchNode().GetIfElse().GetError().Message)
		if err != nil {
			return nil, err
		}
		gb.graphNodes[name] = subNode
		if err := gb.addBranchSubNodeEdge(graph, parentBranchNode, subNode, "orElse - Fail"); err != nil {
			return nil, err
		}
	} else {
		subNode, err := gb.constructNode(prefix, graph, n.GetBranchNode().GetIfElse().GetElseNode())
		if err != nil {
			return nil, err
		}
		if err := gb.addBranchSubNodeEdge(graph, parentBranchNode, subNode, "orElse"); err != nil {
			return nil, err
		}
	}

	if n.GetBranchNode().GetIfElse().GetOther() != nil {
		for _, c := range n.GetBranchNode().GetIfElse().GetOther() {
			subNode, err := gb.constructNode(prefix, graph, c.ThenNode)
			if err != nil {
				return nil, err
			}
			if err := gb.addBranchSubNodeEdge(graph, parentBranchNode, subNode, booleanExprToString(c.Condition)); err != nil {
				return nil, err
			}
		}
	}
	return parentBranchNode, nil
}

func (gb *graphBuilder) constructNode(prefix string, graph *cgraph.Graph, n *core.Node) (*cgraph.Node, error) {
	name := getName(prefix, n.Id)
	var err error
	var gn *cgraph.Node

	if n.Id == "start-node" {
		gn, err = constructStartNode(graph)
	} else if n.Id == "end-node" {
		gn, err = constructEndNode(graph)
	} else {
		switch n.Target.(type) {
		case *core.Node_TaskNode:
			tID := n.GetTaskNode().GetReferenceId().String()
			t, ok := gb.tasks[tID]
			if !ok {
				return nil, fmt.Errorf("failed to find task [%s] in closure", tID)
			}
			gn, err = constructTaskNode(name, graph, n, t)
		case *core.Node_BranchNode:
			branch := graph.SubGraph(fmt.Sprintf("cluster_"+n.Metadata.Name), 2)
			gn, err = gb.constructBranchNode(prefix, branch, n)
			gb.nodeClusters[name] = branch
		case *core.Node_WorkflowNode:
			gn, err = graph.CreateNode(name)
		}
	}
	if err != nil {
		return nil, err
	}
	gb.graphNodes[name] = gn
	return gn, nil
}

func (gb *graphBuilder) addEdge(fromNodeName, toNodeName string, graph *cgraph.Graph) error {
	toNode, toOk := gb.graphNodes[toNodeName]
	fromNode, fromOk := gb.graphNodes[fromNodeName]
	if !toOk || !fromOk {
		return fmt.Errorf("nodes[%s] -> [%s] referenced before creation", fromNodeName, toNodeName)
	}
	edgeName := fmt.Sprintf("%s-%s", fromNodeName, toNodeName)
	if _, ok := gb.graphEdges[edgeName]; !ok {
		edge, err := graph.CreateEdge(edgeName, fromNode, toNode)
		if err != nil {
			return err
		}
		// Now lets check that the toNode or the fromNode is a cluster. If so then following this thread,
		// https://stackoverflow.com/questions/2012036/graphviz-how-to-connect-subgraphs, we will connect the cluster
		if c, ok := gb.nodeClusters[toNodeName]; ok {
			edge.SetLogicalHead(c.Name())
		}
		if c, ok := gb.nodeClusters[fromNodeName]; ok {
			edge.SetLogicalTail(c.Name())
		}
		gb.graphEdges[edgeName] = edge
	}
	return nil
}

func (gb *graphBuilder) constructGraph(prefix string, graph *cgraph.Graph, w *core.CompiledWorkflow) error {
	if w == nil || w.Template == nil {
		return nil
	}
	for _, n := range w.Template.Nodes {
		if _, err := gb.constructNode(prefix, graph, n); err != nil {
			return err
		}
	}

	for name, _ := range gb.graphNodes {
		upstreamNodes := w.Connections.Upstream[name]
		downstreamNodes := w.Connections.Downstream[name]
		if downstreamNodes != nil {
			for _, n := range downstreamNodes.Ids {
				if err := gb.addEdge(name, n, graph); err != nil {
					return err
				}
			}
		}
		if upstreamNodes != nil {
			for _, n := range upstreamNodes.Ids {
				if err := gb.addEdge(n, name, graph); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (gb *graphBuilder) CompiledWorkflowClosureToGraph(g *graphviz.Graphviz, w *core.CompiledWorkflowClosure) (*cgraph.Graph, error) {
	graph, err := g.Graph(graphviz.Directed)
	if err != nil {
		return nil, errors.Wrapf("GraphInitFailure", err, "failed to initialize graphviz")
	}

	graph.SetCompound(true)
	tLookup := make(map[string]*core.CompiledTask)
	for _, t := range w.Tasks {
		tLookup[t.Template.Id.String()] = t
	}
	gb.tasks = tLookup

	return graph, gb.constructGraph("", graph, w.Primary)
}

func newGraphBuilder() *graphBuilder {
	return &graphBuilder{
		graphNodes: make(map[string]*cgraph.Node),
		graphEdges: make(map[string]*cgraph.Edge),
		subWf:      make(map[string]*cgraph.Graph),
		nodeClusters: make(map[string]*cgraph.Graph),
	}
}

// RenderWorkflow Renders the workflow graph to the given file
func RenderWorkflow(w *core.CompiledWorkflowClosure, o graphviz.Format) ([]byte, error) {
	g := graphviz.New()
	defer func() {
		if err := g.Close(); err != nil {
			logger.Fatalf(context.TODO(), "failed to close graphviz. err: %s", err)
		}
	}()
	gb := newGraphBuilder()
	graph, err := gb.CompiledWorkflowClosureToGraph(g, w)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			logger.Fatalf(context.TODO(), "Failed to close the graphviz Graph. err: %s", err)
		}
	}()

	var buf bytes.Buffer
	if err := g.Render(graph, o, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
