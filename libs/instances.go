package libs

import (	
	"github.com/james-nesbitt/coach-tools/log"
)

type InstancesSettings interface {
	UnSerialize() interface{}
}

type Instances interface {
	Init(logger log.Log, client Client, settings InstancesSettings) bool
	Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool

	Instance(name string) (Instance, bool)
}

type BaseInstances struct {
	log log.Log
	instanceMap map[string]Instance
	client Client
}
func (instances *BaseInstances) Init(logger log.Log, client Client, settings InstancesSettings) bool {
	instances.log = logger
	instances.client = client
	instances.instanceMap = map[string]Instance{}
	return true
}
func (instances *BaseInstances) Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool {
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Base Instances")
	return true
}

func (instances *BaseInstances) Instance(id string) (instance Instance, ok bool) {
	instance, ok = instances.instanceMap[id]
	return
}
func (instances *BaseInstances) MakeInstance(id string, machineName string) (Instance, bool) {
	instance := BaseInstance{id: id, machineName: machineName, log: instances.log.MakeChild(id)}
	instance.client = instances.client.InstanceClient(&instance)
	return Instance(&instance), true
}
