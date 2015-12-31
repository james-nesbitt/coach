package libs

import (
	"github.com/james-nesbitt/coach-tools/log"
)

type ClientSettings interface {
	Settings() interface{}
}

type Client interface {
	Init(logger log.Log, settings ClientSettings) bool
	Prepare(logger log.Log, nodes *Nodes, node Node) bool

	NodeClient(node Node) NodeClient
	InstanceClient(instance Instance) InstanceClient
}

/**
 * NodeClient gives a configured client ready to handle
 * client actions for a Node, without further configuration
 * The NodeClient is also used to generate InstanceClients
 * when needed.
 */
type NodeClient interface {
	HasImage() bool // Has this Node got an built or pulled image?

	Info() bool

	Build() bool
	Destroy() bool
	Pull() bool
}

/*
 * IntancsClient gives a configured client ready to handle
 * client actions for an Instance, without further configuration
 */
type InstanceClient interface {
	HasContainer() bool // Does this instance have a matching container
	IsRunning() bool    // Is this instance container running

	Attach(force bool) bool
	Create() bool
	Commit() bool
	Remove(force bool) bool
	Start() bool
	Stop(force bool) bool
	Pause() bool
	Unpause() bool
	Run(overrideCmd []string) bool
}
