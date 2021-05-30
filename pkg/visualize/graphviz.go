package visualize

import (
	"bytes"
	"context"
	"fmt"
	"github.com/flyteorg/flyteidl/gen/pb-go/flyteidl/core"
	"github.com/flyteorg/flytestdlib/errors"
	"github.com/flyteorg/flytestdlib/logger"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

// RenderWorkflow Renders the workflow graph to the given file
func RenderWorkflow(w *core.CompiledWorkflowClosure, file string) error {
	logger.Infof(context.TODO(), "outputting file: %s", file)

	g := graphviz.New()
	graph, err := g.Graph(graphviz.Directed)
	if err != nil {
		return errors.Wrapf("GraphInitFailure", err, "failed to initialize graphviz")
	}
	graph.SetRankDir(cgraph.LRRank)

	defer func() {
		if err = graph.Close(); err != nil {
			logger.Fatalf(context.TODO(), "Failed to close the graphviz Graph. err: %s", err)
		}
	}()

	if err = createNodesAndEdgesFromWorkflow(graph, w) ; err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = g.Render(graph, "dot", &buf); err != nil {
		logger.Fatal(context.TODO(), err)
	}
	fmt.Println(buf.String())

	// 1. write encoded PNG data to buffer
	if err = g.Render(graph, graphviz.SVG, &buf); err != nil {
		logger.Fatal(context.TODO(), err)
	}

	// 3. write to file directly
	if err = g.RenderFilename(graph, graphviz.SVG, file); err != nil {
		logger.Fatal(context.TODO(), err)
	}
	return nil
}

func createNodesAndEdgesFromWorkflow(graph *cgraph.Graph,  w *core.CompiledWorkflowClosure) error {
	graphNodes := make(map[string]*cgraph.Node)
	for _, n := range w.Primary.Template.Nodes {
		name := n.Id
		if n.Id == "start-node" || n.Id == "end-node" {
			node, err := graph.CreateNode(name)
			if err != nil {
				return err
			}
			node.SetShape(cgraph.DoubleCircleShape)
			node.SetColor("#ff00ff")
			graphNodes[name] = node
		} else {
			switch n.Target.(type) {
			case *core.Node_TaskNode:
				node, err := graph.CreateNode(name)
				if err != nil {
					return err
				}
				graphNodes[name] = node
				node.SetColor("#0000ff")
				node.SetShape(cgraph.HouseShape)
			case *core.Node_BranchNode:
				node, err := graph.CreateNode(name)
				if err != nil {
					return err
				}
				node.SetShape(cgraph.DiamondShape)
				node.SetColor("#00ff00")
				graphNodes[name] = node
			case *core.Node_WorkflowNode:
				node, err := graph.CreateNode(name)
				if err != nil {
					return err
				}
				graphNodes[name] = node
				node.SetColor("#ff0000")
				node.SetShape(cgraph.DoubleOctagonShape)
			}
		}
	}
	return createEdges(graphNodes, graph, w)
}

func createEdges(graphNodes map[string]*cgraph.Node, graph *cgraph.Graph, w *core.CompiledWorkflowClosure) error {
	for graphNodeName, graphNode := range graphNodes {
		upstreamNodes, _ := w.Primary.Connections.Upstream[graphNodeName]
		downstreamNodes, _ := w.Primary.Connections.Downstream[graphNodeName]
		if downstreamNodes != nil {
			for _, downStreamNodeId := range downstreamNodes.Ids {
				dNode, ok := graphNodes[downStreamNodeId]
				if !ok {
					return fmt.Errorf("node[%s], downstream from[%s] referenced before creation", downStreamNodeId, graphNodeName)
				}
				edge, err := graph.CreateEdge(fmt.Sprintf("%s-%s", graphNodeName, downStreamNodeId), graphNode, dNode)
				if err != nil {
					return err
				}
				edge.SetLabel(fmt.Sprintf("%s-%s", graphNodeName, downStreamNodeId))
			}
		}
		if upstreamNodes != nil {
			for _, upStreamNodeId := range upstreamNodes.Ids {
				uNode, ok := graphNodes[upStreamNodeId]
				if !ok {
					return fmt.Errorf("node[%s], upstream from[%s] referenced before creation", upStreamNodeId, graphNodeName)
				}
				edge, err := graph.CreateEdge(fmt.Sprintf("%s-%s", upStreamNodeId, graphNodeName), graphNode, uNode)
				if err != nil {
					return err
				}
				edge.SetLabel(fmt.Sprintf("%s-%s", graphNodeName, upStreamNodeId))
			}
		}
	}
	return nil
}