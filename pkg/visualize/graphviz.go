package visualize

import (
	"fmt"
	"strings"

	"github.com/flyteorg/flyteidl/clients/go/coreutils"

	graphviz "github.com/awalterschulze/gographviz"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
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

func constructStartNode(parentGraph string, n string, graph *graphviz.Graph) (*graphviz.Node, error) {
	attrs := map[string]string {"shape" : "doublecircle", "color" : "green"}
	attrs["label"] = "start"
	err := graph.AddNode(parentGraph, n, attrs)
	return graph.Nodes.Lookup[n], err
}

func constructEndNode(parentGraph string,n string, graph *graphviz.Graph) (*graphviz.Node, error) {
	attrs := map[string]string {"shape" : "doublecircle", "color" : "red"}
	attrs["label"] = "end"
	err := graph.AddNode(parentGraph, n, attrs)
	return graph.Nodes.Lookup[n], err
}

func constructTaskNode(parentGraph string, name string, graph *graphviz.Graph, n *core.Node, t *core.CompiledTask) (*graphviz.Node, error) {
	attrs := map[string]string {"shape" : "box"}
	if n.Metadata != nil && n.Metadata.Name != "" {
		v := strings.LastIndexAny(n.Metadata.Name, ".")
		attrs["label"] = fmt.Sprintf("\"%s [%s]\"", n.Metadata.Name[v+1:], t.Template.Type)
	}
	tName := strings.ReplaceAll(name, "-", "_")
	err := graph.AddNode(parentGraph, tName, attrs)
	return graph.Nodes.Lookup[tName], err
}

func constructErrorNode(parentGraph string, name string, graph *graphviz.Graph, m string) (*graphviz.Node, error) {
	attrs := map[string]string {"shape" : "box", "color" : "red", "label" : m}
	eName := strings.ReplaceAll(name, "-", "_")
	err := graph.AddNode(parentGraph, eName, attrs)
	return graph.Nodes.Lookup[eName], err
}

func constructBranchConditionNode(parentGraph string, name string, graph *graphviz.Graph, n *core.Node) (*graphviz.Node, error) {
	attrs := map[string]string {"shape" : "diamond"}
	if n.Metadata != nil && n.Metadata.Name != "" {
		attrs["label"] = n.Metadata.Name
	}
	cName := strings.ReplaceAll(name, "-", "_")
	err := graph.AddNode(parentGraph, cName, attrs)
	return graph.Nodes.Lookup[cName], err
}

func getName(prefix, id string) string {
	if prefix != "" {
		return prefix + "_" + id
	}
	return id
}

type graphBuilder struct {
	// Mutated as graph is built
	graphNodes map[string]*graphviz.Node
	// Mutated as graph is built. lookup table for all graphviz compiled edges.
	graphEdges map[string]*graphviz.Edge
	// lookup table for all graphviz compiled subgraphs
	subWf map[string]*core.CompiledWorkflow
	// a lookup table for all tasks in the graph
	tasks map[string]*core.CompiledTask
}

func (gb *graphBuilder) addBranchSubNodeEdge(graph *graphviz.Graph, parentNode, n *graphviz.Node, label string) error {
	edgeName := fmt.Sprintf("%s-%s", parentNode.Name, n.Name)
	if _, ok := gb.graphEdges[edgeName]; !ok {
		attrs := map[string]string {}
		attrs["label"] = fmt.Sprintf("\"%s\"", label)
		err := graph.AddEdge(parentNode.Name, n.Name, true, attrs)
		if err != nil {
			return err
		}
		gb.graphEdges[edgeName] = graph.Edges.SrcToDsts[parentNode.Name][n.Name][0]
	}
	return nil
}

func (gb *graphBuilder) constructBranchNode(parentGraph string, prefix string, graph *graphviz.Graph, n *core.Node) (*graphviz.Node, error) {
	parentBranchNode, err := constructBranchConditionNode(parentGraph, getName(prefix, n.Id), graph, n)
	if err != nil {
		return nil, err
	}
	gb.graphNodes[n.Id] = parentBranchNode

	if n.GetBranchNode().GetIfElse() == nil {
		return parentBranchNode, nil
	}

	subNode, err := gb.constructNode(parentGraph, prefix, graph, n.GetBranchNode().GetIfElse().Case.ThenNode)
	if err != nil {
		return nil, err
	}
	if err := gb.addBranchSubNodeEdge(graph, parentBranchNode, subNode, booleanExprToString(n.GetBranchNode().GetIfElse().Case.Condition)); err != nil {
		return nil, err
	}

	if n.GetBranchNode().GetIfElse().GetError() != nil {
		name := fmt.Sprintf("%s-error", parentBranchNode.Name)
		subNode, err := constructErrorNode(prefix, name, graph, n.GetBranchNode().GetIfElse().GetError().Message)
		if err != nil {
			return nil, err
		}
		gb.graphNodes[name] = subNode
		if err := gb.addBranchSubNodeEdge(graph, parentBranchNode, subNode, "orElse - Fail"); err != nil {
			return nil, err
		}
	} else {
		subNode, err := gb.constructNode(parentGraph, prefix, graph, n.GetBranchNode().GetIfElse().GetElseNode())
		if err != nil {
			return nil, err
		}
		if err := gb.addBranchSubNodeEdge(graph, parentBranchNode, subNode, "orElse"); err != nil {
			return nil, err
		}
	}

	if n.GetBranchNode().GetIfElse().GetOther() != nil {
		for _, c := range n.GetBranchNode().GetIfElse().GetOther() {
			subNode, err := gb.constructNode(parentGraph, prefix, graph, c.ThenNode)
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

func (gb *graphBuilder) constructNode(parentGraphName string, prefix string, graph *graphviz.Graph, n *core.Node) (*graphviz.Node, error) {
	name := getName(prefix, n.Id)
	var err error
	var gn *graphviz.Node

	if n.Id == "start-node" {
		gn, err = constructStartNode(parentGraphName, strings.ReplaceAll(name, "-", "_"), graph)
	} else if n.Id == "end-node" {
		gn, err = constructEndNode(parentGraphName, strings.ReplaceAll(name, "-", "_"), graph)
	} else {
		switch n.Target.(type) {
		case *core.Node_TaskNode:
			tID := n.GetTaskNode().GetReferenceId().String()
			t, ok := gb.tasks[tID]
			if !ok {
				return nil, fmt.Errorf("failed to find task [%s] in closure", tID)
			}
			gn, err = constructTaskNode(parentGraphName, name, graph, n, t)
			if err != nil {
				return nil, err
			}
		case *core.Node_BranchNode:
			err := graph.AddSubGraph(parentGraphName, fmt.Sprintf("cluster_"+n.Metadata.Name), nil)
			if err != nil {
				return nil, err
			}
			gn, err = gb.constructBranchNode(fmt.Sprintf("cluster_"+n.Metadata.Name), prefix, graph, n)
			if err != nil {
				return nil, err
			}
		case *core.Node_WorkflowNode:
			if n.GetWorkflowNode().GetLaunchplanRef() != nil {
				attrs := map[string]string {}
				err := graph.AddNode(parentGraphName, name, attrs)
				if err != nil {
					return nil, err
				}
			} else {
				err := graph.AddSubGraph(parentGraphName, "cluster_"+name, nil)
				if err != nil {
					return nil, err
				}
				subGB := graphBuilderFromParent(gb)
				swf, ok := gb.subWf[n.GetWorkflowNode().GetSubWorkflowRef().String()]
				if !ok {
					return nil, fmt.Errorf("subworkfow [%s] not found", n.GetWorkflowNode().GetSubWorkflowRef().String())
				}
				if err := subGB.constructGraph("cluster_"+name, name, graph, swf); err != nil {
					return nil, err
				}
				gn = subGB.graphNodes["start-node"]
			}
		}
	}
	if err != nil {
		return nil, err
	}
	gb.graphNodes[n.Id] = gn
	return gn, nil
}

func (gb *graphBuilder) addEdge(fromNodeName, toNodeName string, graph *graphviz.Graph) error {
	toNode, toOk := gb.graphNodes[toNodeName]
	fromNode, fromOk := gb.graphNodes[fromNodeName]
	if !toOk || !fromOk {
		return fmt.Errorf("nodes[%s] -> [%s] referenced before creation", fromNodeName, toNodeName)
	}
	if _,ok := graph.Edges.SrcToDsts[fromNode.Name][toNode.Name]; !ok {
		attrs := map[string]string {}
		err := graph.AddEdge(fromNode.Name, toNode.Name, true, attrs)
		if err != nil {
			return err
		}
		// Now lets check that the toNode or the fromNode is a cluster. If so then following this thread,
		// https://stackoverflow.com/questions/2012036/graphviz-how-to-connect-subgraphs, we will connect the cluster

		if c, ok := graph.Nodes.Lookup[fromNode.Name]; ok {
			attrs["ltail"] = c.Name
		}
		if c, ok := graph.Nodes.Lookup[toNode.Name]; ok {
			attrs["lhead"] = c.Name
		}
	}
	return nil
}

func (gb *graphBuilder) constructGraph(parentGraphName string, prefix string, graph *graphviz.Graph, w *core.CompiledWorkflow) error {
	if w == nil || w.Template == nil {
		return nil
	}
	for _, n := range w.Template.Nodes {
		if _, err := gb.constructNode(parentGraphName, prefix, graph, n); err != nil {
			return err
		}
	}

	for name := range gb.graphNodes {
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

func (gb *graphBuilder) CompiledWorkflowClosureToGraph(w *core.CompiledWorkflowClosure) (*graphviz.Graph, error) {
	dotGraph := graphviz.NewGraph()
	_ = dotGraph.SetDir(true)
	_ = dotGraph.SetStrict(true)

	tLookup := make(map[string]*core.CompiledTask)
	for _, t := range w.Tasks {
		tLookup[t.Template.Id.String()] = t
	}
	gb.tasks = tLookup
	wLookup := make(map[string]*core.CompiledWorkflow)
	for _, swf := range w.SubWorkflows {
		wLookup[swf.Template.Id.String()] = swf
	}
	gb.subWf = wLookup

	return dotGraph, gb.constructGraph("", "", dotGraph, w.Primary)
}

func newGraphBuilder() *graphBuilder {
	return &graphBuilder{
		graphNodes:   make(map[string]*graphviz.Node),
		graphEdges:   make(map[string]*graphviz.Edge),
	}
}

func graphBuilderFromParent(gb *graphBuilder) *graphBuilder {
	newGB := newGraphBuilder()
	newGB.subWf = gb.subWf
	newGB.tasks = gb.tasks
	return newGB
}

// RenderWorkflow Renders the workflow graph to the given file
func RenderWorkflow(w *core.CompiledWorkflowClosure) ([]byte, error) {
	gb := newGraphBuilder()
	graph, err := gb.CompiledWorkflowClosureToGraph(w)
	if err != nil {
		return nil, err
	}

	fmt.Printf(graph.String())
	return []byte(graph.String()), nil
}
