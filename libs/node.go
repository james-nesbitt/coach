package libs

import (
	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

type Node interface {
	Type() string
	Id() string
	MachineName() string

	Init(logger log.Log, name string, project *conf.Project, client Client, instancesSettings InstancesSettings) bool
	Prepare(logger log.Log, nodes *Nodes) bool

	Can(action string) bool

	AddDependency(target string) bool
	DependsOn(target string) bool

	Client() NodeClient
	Instances() Instances
}

type BaseNode struct {
	log                log.Log
	conf               *conf.Project
	name               string
	machineName        string
	client             Client
	instances          Instances
	manualDependencies []string
}

// Declare node type
func (node *BaseNode) Type() string {
	return "base"
}
func (node *BaseNode) Id() string {
	return node.name
}
func (node *BaseNode) MachineName() string {
	return node.conf.Name + "_" + node.name
}

// Constructor for BaseNode
func (node *BaseNode) Init(logger log.Log, name string, project *conf.Project, client Client, instancesSettings InstancesSettings) bool {
	node.log = logger
	node.conf = project
	node.name = name
	node.client = client
	node.manualDependencies = []string{}

	instancesMachineName := node.MachineName()

	settingsInterface := instancesSettings.Settings()
	switch settingsInterface.(type) {
	case FixedInstancesSettings:
		node.instances = Instances(&FixedInstances{})
		instancesMachineName += "_fixed_"
	case TemporaryInstancesSettings:
		node.instances = Instances(&TemporaryInstances{})
		instancesMachineName += "_temp_"
	case ScaledInstancesSettings:
		node.instances = Instances(&ScaledInstances{})
		instancesMachineName += "_scaled_"
	case SingleInstancesSettings:
		node.instances = Instances(&SingleInstances{})
	default:
		node.instances = Instances(&SingleInstances{})
	}

	node.instances.Init(logger, instancesMachineName, client, instancesSettings)

	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Built new node:", node.client)
	return true
}

// Post initialization preparation
func (node *BaseNode) Prepare(logger log.Log, nodes *Nodes) (success bool) {
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Base Node")
	success = true

	node.log = logger

	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Preparing Client", nil)
	success = success && node.client.Prepare(logger.MakeChild("client"), nodes, node)

	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Preparing Instances", nil)
	success = success && node.instances.Prepare(logger.MakeChild("instances"), node.client, nodes, node)

	return success
}

func (node *BaseNode) Can(action string) bool {
	switch action {
	case "run":
		return false
	default:
		return node.client.Can(action)
	}
}

func (node *BaseNode) Client() NodeClient {
	if node.client == nil {
		return nil
	}
	return node.client.NodeClient(node)
}

func (node *BaseNode) Instances() Instances {
	return node.instances
}

func (node *BaseNode) AddDependency(target string) bool {
	if !node.hasManualDependency(target) {
		node.manualDependencies = append(node.manualDependencies, target)
	}
	return true
}
func (node *BaseNode) DependsOn(target string) bool {
	return node.hasManualDependency(target) || node.client.DependsOn(target)
}
func (node *BaseNode) hasManualDependency(target string) bool {
	for _, manualDependency := range node.manualDependencies {
		if manualDependency == target {
			return true
		}
	}
	return false
}
