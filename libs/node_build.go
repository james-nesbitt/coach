package libs

import (
	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

type BuildNodeSettings struct {
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
}

type BuildNode struct {
	BaseNode
	Settings BuildNodeSettings
}

// Declare node type
func (node *BuildNode) Type() string {
	return "build"
}
func (node *BuildNode) Init(logger log.Log, name string, project *conf.Project, client Client, instancesSettings InstancesSettings) bool {
	node.BaseNode.Init(logger, name, project, client, instancesSettings)

	settingsInterface := instancesSettings.Settings()
	switch settingsInterface.(type) {
	case FixedInstancesSettings:
		logger.Warning("Build node cannot be configured to use fixed instances.  Using null instance instead.")
		node.defaultInstances(logger, client, instancesSettings)
	case ScaledInstancesSettings:
		logger.Warning("Build node cannot be configured to use scaled instances.  Using null instance instead.")
		node.defaultInstances(logger, client, instancesSettings)
	case SingleInstancesSettings:
		logger.Warning("Build node cannot be configured to use single instances.  Using null instance instead.")
		node.defaultInstances(logger, client, instancesSettings)
	case TemporaryInstancesSettings:
		logger.Warning("Build node cannot be configured to use disposable instances.  Using null instance instead.")
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
func (node *BuildNode) defaultInstances(logger log.Log, client Client, instancesSettings InstancesSettings) {
	node.instances = Instances(&NullInstances{})
}

// Build Nodes can only build
func (node *BuildNode) Can(action string) bool {
	switch action {
	case "destroy":
		fallthrough
	case "clean":
		fallthrough
	case "build":
		return true
	default:
		return false
	}
}

// Return some string status for the node
func (node *BuildNode) Status(logger log.Log) []string {
	status := []string{"BUILD"}

	if node.Client().HasImage() {
		status = append(status, "Image:BUILT")
	} else {
		status = append(status, "Image:NOT-BUILT")
	}

	return status
}