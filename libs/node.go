package libs

import (
	"github.com/james-nesbitt/coach-tools/log"
)

type Node interface {
	Init(logger log.Log, client Client, instances Instances) bool
	Prepare(logger log.Log, nodes *Nodes) bool

	Can(action string) bool

	NodeClient() NodeClient

	Instances() Instances
}

type BaseNode struct {
	log log.Log
	client Client
	instances Instances
}
// Constructor for BaseNode
func (node *BaseNode) Init(logger log.Log, client Client, instances Instances) bool {
	node.log = logger
	node.client = client
	node.instances = instances

	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Built new node:", node.client)
	return true
}
// Post initialization preparation
func (node *BaseNode) Prepare(logger log.Log, nodes *Nodes) (success bool)	 {
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Base Node")
	success = true

	node.log = logger

	success = success && node.client.Prepare(logger.MakeChild("client"),nodes, node)
	success = success && node.instances.Prepare(logger.MakeChild("instances"), node.client, nodes, node)

	return success
}

func (node *BaseNode) Can(action string) bool {
	return true
}

func (node *BaseNode) NodeClient() NodeClient {
	if node.client==nil {
		return nil
	}
	return node.client.NodeClient(node)
}

func (node *BaseNode) Instances() Instances {
	return node.instances
}
