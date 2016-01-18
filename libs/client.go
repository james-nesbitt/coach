package libs

import (
	"github.com/james-nesbitt/coach-tools/conf"
	"github.com/james-nesbitt/coach-tools/log"
)

type ClientSettings interface {
	Settings() interface{}
}

type Client interface {
	Init(logger log.Log, project *conf.Project, settings ClientSettings) bool
	Prepare(logger log.Log, nodes *Nodes, node Node) bool

	Can(action string) bool

	NodeClient(node Node) NodeClient
	InstancesClient(instances Instances) InstancesClient
	InstanceClient(instance Instance) InstanceClient

	DependsOn(target string) bool
}

/**
 * NodeClient gives a configured client ready to handle
 * client actions for a Node, without further configuration
 * The NodeClient is also used to generate InstanceClients
 * when needed.
 */
type NodeClient interface {
	Can(action string) bool
	HasImage() bool // Has this Node got an built or pulled image?

	NodeInfo(logger log.Log)

	Build(logger log.Log, force bool) bool
	Destroy(logger log.Log) bool
	Pull(logger log.Log, force bool) bool
}

/**
 *
 */
type InstancesClient interface {
	Can(action string) bool
	InstancesFound(logger log.Log) []string

	InstancesInfo(logger log.Log)
}

/*
 * IntancsClient gives a configured client ready to handle
 * client actions for an Instance, without further configuration
 */
type InstanceClient interface {
	Can(action string) bool
	HasContainer() bool // Does this instance have a matching container
	IsRunning() bool    // Is this instance container running

	Attach(logger log.Log, force bool) bool
	Create(logger log.Log, overrideCmd []string, force bool) bool
	Commit(logger log.Log) bool
	Remove(logger log.Log, force bool) bool
	Start(logger log.Log, force bool) bool
	Stop(logger log.Log, force bool, timeout uint) bool
	Pause(logger log.Log) bool
	Unpause(logger log.Log) bool
	Run(logger log.Log, overrideCmd []string) bool
}
