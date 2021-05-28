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

func constructGraph(prefix string, graph *cgraph.Graph, w *core.CompiledWorkflow, subWf map[string]*cgraph.Graph) error {
	grapNodes := make(map[string]*cgraph.Node)
	for _, n := range w.Template.Nodes {
		name := n.Id
		if prefix != "" {
			name = prefix + "-" + n.Id
		}

		if n.Id == "start-node" || n.Id == "end-node"{
			gn, err := graph.CreateNode(name)
			if err != nil {
				return err
			}
			grapNodes[name] = gn
		} else {
			switch n.Target.(type) {
			case *core.Node_TaskNode:
				gn, err := graph.CreateNode(name)
				if err != nil {
					return err
				}
				grapNodes[name] = gn
			case *core.Node_BranchNode:
				gn, err := graph.CreateNode(name)
				if err != nil {
					return err
				}
				grapNodes[name] = gn
			case *core.Node_WorkflowNode:
				gn, err := graph.CreateNode(name)
				if err != nil {
					return err
				}
				grapNodes[name] = gn
			}
		}
	}

	for name, node := range grapNodes {
		upstreamNodes, _ := w.Connections.Upstream[name]
		downstreamNodes, _ := w.Connections.Downstream[name]
		if downstreamNodes != nil {
			for _, n := range downstreamNodes.Ids {
				dNode, ok := grapNodes[n]
				if !ok {
					return fmt.Errorf("node[%s], downstream from[%s] referenced before creation", n, name)
				}
				_, err := graph.CreateEdge(fmt.Sprintf("%s-%s", name, n), node, dNode)
				if err != nil {
					return err
				}
			}
		}
		if upstreamNodes != nil {
			for _, n := range upstreamNodes.Ids {
				uNode, ok := grapNodes[n]
				if !ok {
					return fmt.Errorf("node[%s], upstream from[%s] referenced before creation", n, name)
				}
				_, err := graph.CreateEdge(fmt.Sprintf("%s-%s", n, name), node, uNode)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func CompiledWorkflowClosureToGraph(g *graphviz.Graphviz, w *core.CompiledWorkflowClosure) (*cgraph.Graph, error) {
	graph, err := g.Graph(graphviz.Directed)
	if err != nil {
		return nil, errors.Wrapf("GraphInitFailure", err, "failed to initialize graphviz")
	}

	defer func() {
		if err := graph.Close(); err != nil {
			logger.Fatalf(context.TODO(), "Failed to close the graphviz Graph. err: %s", err)
		}
	}()

	return graph, constructGraph("", graph, w.Primary, nil)
}

// RenderWorkflow Renders the workflow graph to the given file
func RenderWorkflow(w *core.CompiledWorkflowClosure, file string) error {
	g := graphviz.New()
	defer func() {
		if err := g.Close(); err != nil {
			logger.Fatalf(context.TODO(), "failed to close graphviz. err: %s", err)
		}
	}()
	graph, err := CompiledWorkflowClosureToGraph(g, w)
	if err != nil {
		return err
	}

	logger.Infof(context.TODO(), "outputing file: %s", file)
	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.XDOT, &buf); err != nil {
		return err
	}
	logger.Infof(context.TODO(), buf.String())
	if err := g.RenderFilename(graph, graphviz.SVG, file); err != nil {
		return err
	}
	return nil
}
