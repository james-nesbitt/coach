package libs

import (
	"github.com/james-nesbitt/coach-tools/conf"
	"github.com/james-nesbitt/coach-tools/log"
)

type Nodes struct {
	NodesMap map[string]*Node
}

func MakeNodes(logger log.Log, project *conf.Project, clientFactories *ClientFactories) *Nodes {
	nodes := &Nodes{NodesMap: map[string]*Node{}}

	nodes.from_NodesYaml(logger.MakeChild("yaml"), project, false)

	return nodes
}


// If node exists, return a reference to it.
func (nodes *Nodes) Node(name string) (node *Node, exists bool) {
	node, exists = nodes.NodesMap[name]
	return
}