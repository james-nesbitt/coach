package libs

import (
	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

type ServiceNodeSettings struct {
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
}

type ServiceNode struct {
	Settings ServiceNodeSettings
	BaseNode
}

// Declare node type
func (node *ServiceNode) Type() string {
	return "service"
}
func (node *ServiceNode) Init(logger log.Log, name string, project *conf.Project, client Client, instancesSettings InstancesSettings) bool {
	node.BaseNode.Init(logger, name, project, client, instancesSettings)

	settingsInterface := instancesSettings.Settings()
	switch settingsInterface.(type) {
	case FixedInstancesSettings:
		node.instances = Instances(&FixedInstances{})
	case ScaledInstancesSettings:
		node.instances = Instances(&ScaledInstances{})
	case TemporaryInstancesSettings:
		logger.Warning("Service node cannot be configured to use disposable instances.  Using single instance instead.")
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
func (node *ServiceNode) defaultInstances(logger log.Log, client Client, instancesSettings InstancesSettings) {
	node.instances = Instances(&SingleInstances{})
}

// Build Nodes can only build
func (node *ServiceNode) Can(action string) bool {
	switch action {
	default:
		return node.BaseNode.Can(action)
	}
}

// Return some string status for the node
func (node *ServiceNode) Status(logger log.Log) []string {
	status := []string{"SERVICE"}
	status = append(status, node.BaseNode.Status(logger)...)
	return status
}