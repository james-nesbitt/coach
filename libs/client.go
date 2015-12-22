package libs

type NodeClientSettings interface {
	FindRelatedNodes(nodes Nodes)           // Match and retain any node dependencies that the client thinks it may need
	MakeInstanceSettings(instance Instance) // Make an instance of the client for a Node Instance, from the Node configuration
}

/**
 * NodeClient gives a configured client ready to handle
 * client actions for a Node, without further configuration
 * The NodeClient is also used to generate InstanceClients
 * when needed.
 */
type NodeClient interface {
	TakeNodeSettings(settings NodeClientSettings)

	HasImage() bool // Has this Node got an built or pulled image?

	Info() bool

	Build() bool
	Destroy() bool
	Pull() bool
}

type InstanceClientSettings interface {
	FindRelatedNodes(nodes Nodes) // Match and retain any node dependencies that the client thinks it may need
}

/*
 * IntancsClient gives a configured client ready to handle
 * client actions for an Instance, without further configuration
 */
type InstanceClient interface {
	TakeInstanceSettings(settings InstanceClientSettings)

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
