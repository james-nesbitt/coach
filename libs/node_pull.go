package libs

import (
	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

type PullNodeSettings struct {
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
}

type PullNode struct {
	BaseNode
	Settings PullNodeSettings
}

// Declare node type
func (node *PullNode) Type() string {
	return "pull"
}
func (node *PullNode) Init(logger log.Log, name string, project *conf.Project, client Client, instancesSettings InstancesSettings) bool {
	node.BaseNode.Init(logger, name, project, client, instancesSettings)

	settingsInterface := instancesSettings.Settings()
	switch settingsInterface.(type) {
	case FixedInstancesSettings:
		logger.Warning("Pull node cannot be configured to use fixed instances.  Using null instance instead.")
		node.defaultInstances(logger, client, instancesSettings)
	case ScaledInstancesSettings:
		logger.Warning("Pull node cannot be configured to use scaled instances.  Using null instance instead.")
		node.defaultInstances(logger, client, instancesSettings)
	case SingleInstancesSettings:
		logger.Warning("Pull node cannot be configured to use single instances.  Using null instance instead.")
		node.defaultInstances(logger, client, instancesSettings)
	case TemporaryInstancesSettings:
		logger.Warning("Pull node cannot be configured to use disposable instances.  Using null instance instead.")
		node.defaultInstances(logger, client, instancesSettings)
	default:
		node.defaultInstances(logger, client, instancesSettings)
	}

	node.instances.Init(logger, node.MachineName(), client, instancesSettings)

	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Built new node:", node.client)
	return true
}

// a central instances configuration method, that will create the default single settings
// @note we use this mainly because you cannot fallthrough on a .(type) switch statement
func (node *PullNode) defaultInstances(logger log.Log, client Client, instancesSettings InstancesSettings) {
	node.instances = Instances(&NullInstances{})
}

// Pull Nodes can only Pull
func (node *PullNode) Can(action string) bool {
	switch action {
	case "destroy":
		fallthrough
	case "clean":
		fallthrough
	case "pull":
		return true
	default:
		return false
	}
}
