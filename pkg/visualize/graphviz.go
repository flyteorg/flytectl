package visualize

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/flyteorg/flytestdlib/errors"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func constructStartNode(graph *cgraph.Graph) (*cgraph.Node, error)  {
	gn, err := graph.CreateNode("start-node")
	if err != nil {
		return nil, err
	}
	gn.SetLabel("start")
	gn.SetShape(cgraph.CircleShape)
	gn.SetColor("green")
	return gn, nil
}

func constructEndNode(graph *cgraph.Graph) (*cgraph.Node, error)  {
	gn, err := graph.CreateNode("end-node")
	if err != nil {
		return nil, err
	}
	gn.SetLabel("end")
	gn.SetShape(cgraph.CircleShape)
	gn.SetColor("red")
	return gn, nil
}

func constructTaskNode(name string, graph *cgraph.Graph, n *core.Node, t *core.CompiledTask) (*cgraph.Node, error)  {
	// TODO Add task task type and other information from the task template
	gn, err := graph.CreateNode(name)
	if err != nil {
		return nil, err
	}
	if n.Metadata != nil && n.Metadata.Name != ""{
		v := strings.LastIndexAny(n.Metadata.Name, ".")
		gn.SetLabel(n.Metadata.Name[v+1:])
	}
	gn.SetShape(cgraph.BoxShape)
	return gn, nil
}

func constructErrorNode(name string, graph *cgraph.Graph, m string) (*cgraph.Node, error)  {
	gn, err := graph.CreateNode(name)
	if err != nil {
		return nil, err
	}
	gn.SetLabel(m)
	gn.SetShape(cgraph.PentagonShape)
	return gn, nil
}

func constructBranchConditionNode(name string, graph *cgraph.Graph, n *core.Node) (*cgraph.Node, error)  {
	gn, err := graph.CreateNode(name)
	if err != nil {
		return nil, err
	}
	if n.Metadata != nil && n.Metadata.Name != ""{
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
	graphNodes map[string]*cgraph.Node
	graphEdges map[string]*cgraph.Edge
	subWf      map[string]*cgraph.Graph
	// TODO Add tasks
}

func (gb *graphBuilder) addSubNodeEdge(graph *cgraph.Graph, parentNode, n *cgraph.Node) error {
	edgeName := fmt.Sprintf("%s-%s", parentNode.Name(), n.Name())
	if _, ok := gb.graphEdges[edgeName]; !ok {
		edge, err := graph.CreateEdge(edgeName, parentNode, n)
		if err != nil {
			return err
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
	if err := gb.addSubNodeEdge(graph, parentBranchNode, subNode); err != nil {
		return nil, err
	}

	if n.GetBranchNode().GetIfElse().GetError() != nil {
		name := fmt.Sprintf("%s-error", parentBranchNode.Name())
		subNode, err := constructErrorNode(name, graph, n.GetBranchNode().GetIfElse().GetError().Message)
		if err != nil {
			return nil, err
		}
		gb.graphNodes[name] = subNode
		if err := gb.addSubNodeEdge(graph, parentBranchNode, subNode); err != nil {
			return nil, err
		}
	} else {
		subNode, err := gb.constructNode(prefix, graph, n.GetBranchNode().GetIfElse().GetElseNode())
		if err != nil {
			return nil, err
		}
		if err := gb.addSubNodeEdge(graph, parentBranchNode, subNode); err != nil {
			return nil, err
		}
	}

	if n.GetBranchNode().GetIfElse().GetOther() != nil {
		for _, c := range n.GetBranchNode().GetIfElse().GetOther() {
			subNode, err := gb.constructNode(prefix, graph, c.ThenNode)
			if err != nil {
				return nil, err
			}
			if err := gb.addSubNodeEdge(graph, parentBranchNode, subNode); err != nil {
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

	if n.Id == "start-node"{
		gn, err = constructStartNode(graph)
	} else if n.Id == "end-node" {
		gn, err = constructEndNode(graph)
	} else {
		switch n.Target.(type) {
		case *core.Node_TaskNode:
			// TODO add task lookup
			gn, err = constructTaskNode(name, graph, n, nil)
		case *core.Node_BranchNode:
			branch := graph.SubGraph(fmt.Sprintf("cluster_"+n.Metadata.Name), 2)
			gn, err = gb.constructBranchNode(prefix, branch, n)
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

func (gb *graphBuilder) constructGraph(prefix string, graph *cgraph.Graph, w *core.CompiledWorkflow) error {
	if w == nil || w.Template == nil {
		return nil
	}
	for _, n := range w.Template.Nodes {
		if _, err := gb.constructNode(prefix, graph, n); err != nil {
			return err
		}
	}

	for name, node := range gb.graphNodes {
		upstreamNodes := w.Connections.Upstream[name]
		downstreamNodes := w.Connections.Downstream[name]
		if downstreamNodes != nil {
			for _, n := range downstreamNodes.Ids {
				dNode, ok := gb.graphNodes[n]
				if !ok {
					return fmt.Errorf("node[%s], downstream from[%s] referenced before creation", n, name)
				}
				edgeName := fmt.Sprintf("%s-%s", name, n)
				if _, ok := gb.graphEdges[edgeName]; !ok {
					edge, err := graph.CreateEdge(edgeName, node, dNode)
					if err != nil {
						return err
					}
					gb.graphEdges[edgeName] = edge
				}
			}
		}
		if upstreamNodes != nil {
			for _, n := range upstreamNodes.Ids {
				uNode, ok := gb.graphNodes[n]
				if !ok {
					return fmt.Errorf("node[%s], upstream from[%s] referenced before creation", n, name)
				}
				edgeName := fmt.Sprintf("%s-%s", n, name)
				if _, ok := gb.graphEdges[edgeName]; !ok {
					edge, err := graph.CreateEdge(edgeName, uNode, node)
					if err != nil {
						return err
					}
					gb.graphEdges[edgeName] = edge
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

	return graph, gb.constructGraph("", graph, w.Primary)
}

func newGraphBuilder() *graphBuilder {
	return &graphBuilder{
		graphNodes: make(map[string]*cgraph.Node),
		graphEdges: make(map[string]*cgraph.Edge),
		subWf:      make(map[string]*cgraph.Graph),
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
