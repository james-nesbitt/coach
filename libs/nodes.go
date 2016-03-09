package libs

import (
	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

/**
 * Nodes factory method
 */
func MakeNodes(logger log.Log, project *conf.Project, clientFactories *ClientFactories) *Nodes {
	nodes := &Nodes{}
	nodes.Init(logger)

	nodes.from_NodesYaml(logger.MakeChild("yaml"), project, clientFactories, true)

	return nodes
}

/**
 * NODES: A collection of Node objects considered a set
 */
type Nodes struct {
	NodesMap   map[string]Node
	NodesOrder []string
}

// Initialize a Nodes list
func (nodes *Nodes) Init(logger log.Log) bool {
	nodes.NodesMap = map[string]Node{}
	return true
}
func (nodes *Nodes) Prepare(logger log.Log) bool {
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Nodes")
	for _, id := range nodes.NodesOrder {
		node := nodes.NodesMap[id]
		node.Prepare(logger.MakeChild(id), nodes)
	}

	return true
}

// If node exists, return it.
func (nodes *Nodes) Node(name string) (node Node, exists bool) {
	node, exists = nodes.NodesMap[name]
	return
}

// Attach a Node to the nodes list
func (nodes *Nodes) SetNode(name string, node Node, overwrite bool) bool {
	if _, exists := nodes.NodesMap[name]; exists && !overwrite {
		return false
	}
	nodes.NodesOrder = append(nodes.NodesOrder, name)
	nodes.NodesMap[name] = node
	return true
}
// Disable a Node in the nodes list
func (nodes *Nodes) DisableNode(name string) bool {
	if _, exists := nodes.NodesMap[name]; !exists {
		nodes.log.Warning("Tried to remove node ["+name+"] but it was not found")
		return false
	}
	
	for index, orderName := range nodes.NodesOrder {
		if name==orderName {
			orderLen := len(nodes.NodesOrder)

			if orderLen==1 {
				nodes.NodesOrder = []string{}
			} else if index==0 {
				nodes.NodesOrder = nodes.NodesOrder[1:]
			} else if orderLen==index+1 {
				nodes.NodesOrder = nodes.NodesOrder[:index]
			} else {
				nodes.NodesOrder = append(nodes.NodesOrder[:index], nodes.NodesOrder[index+1:]...)
			}

			delete(nodes.NodesMap, name)
			return true
		}
	}
	nodes.log.Warning("Tried to remove node ["+name+"] but it was not found in the nodes order")
	return true
}

// return a string slice of node keys
func (nodes *Nodes) NodeNames() []string {
	return nodes.NodesOrder
}

// return a rangeable set of Node objects
func (nodes *Nodes) Nodes() []Node {
	ordered := []Node{}
	for _, node := range nodes.NodesOrder {
		ordered = append(ordered, nodes.NodesMap[node])
	}
	return ordered
}
