package libs

import (
	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

type CommandNodeSettings struct {
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
}

type CommandNode struct {
	Settings CommandNodeSettings
	BaseNode
}

// Declare node type
func (node *CommandNode) Type() string {
	return "command"
}
func (node *CommandNode) Init(logger log.Log, name string, project *conf.Project, client Client, instancesSettings InstancesSettings) bool {
	node.BaseNode.Init(logger, name, project, client, instancesSettings)

	settingsInterface := instancesSettings.Settings()
	switch settingsInterface.(type) {
	case FixedInstancesSettings:
		node.instances = Instances(&FixedInstances{})
	case ScaledInstancesSettings:
		logger.Warning("Command node cannot be configured to use disposable instances.  Using single instance instead.")
		node.defaultInstances(logger, client, instancesSettings)
	case SingleInstancesSettings:
		node.instances = Instances(&SingleInstances{})
	default:
		node.defaultInstances(logger, client, instancesSettings)
	}

	node.instances.Init(logger, node.MachineName(), client, instancesSettings)

	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Built new node:", node.client)
	return true
}

// a central instances configuration method, that will create the default single settings
// @note we use this mainly because you cannot fallthrough on a .(type) switch statement
func (node *CommandNode) defaultInstances(logger log.Log, client Client, instancesSettings InstancesSettings) {
	node.instances = Instances(&TemporaryInstances{})
}

func (node *CommandNode) Can(action string) bool {
	switch action {
	case "create":
		fallthrough
	case "start":
		return false
	case "run":
		return true
	default:
		return node.client.Can(action)
	}
}
