package libs

import (
	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

type VolumeNodeSettings struct {
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
}

type VolumeNode struct {
	Settings VolumeNodeSettings
	BaseNode
}

// Declare node type
func (node *VolumeNode) Type() string {
	return "volume"
}
func (node *VolumeNode) Init(logger log.Log, name string, project *conf.Project, client Client, instancesSettings InstancesSettings) bool {
	node.BaseNode.Init(logger, name, project, client, instancesSettings)

	settingsInterface := instancesSettings.Settings()
	switch settingsInterface.(type) {
	case FixedInstancesSettings:
		node.instances = Instances(&FixedInstances{})
	case ScaledInstancesSettings:
		node.instances = Instances(&ScaledInstances{})
	case SingleInstancesSettings:
		node.defaultInstances(logger, client, instancesSettings)
	case TemporaryInstancesSettings:
		logger.Warning("Volume node cannot be configured to use disposable instances.  Using single instance instead.")
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
func (node *VolumeNode) defaultInstances(logger log.Log, client Client, instancesSettings InstancesSettings) {
	node.instances = Instances(&SingleInstances{})
}

// Volume Nodes can only build and create, they can never start
func (node *VolumeNode) Can(action string) bool {
	switch action {
	case "run":
		fallthrough
	case "unpause":
		fallthrough
	case "pause":
		fallthrough
	case "stop":
		fallthrough
	case "start":
		return false
	default:
		return node.BaseNode.Can(action)
	}
}
