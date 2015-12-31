package libs

import (
	"github.com/james-nesbitt/coach-tools/log"
)

const (
	INSTANCE_SINGLE_ID = "single"
	INSTANCE_SINGLE_MACHINENAME = ""
)

type SingleInstancesSettings struct {
	Name string 
}
func (settings SingleInstancesSettings) UnSerialize() interface{} {
	return settings
}

type SingleInstances struct {
	log log.Log
	settings SingleInstancesSettings
	instance SingleInstance
	client Client
}
func (instances *SingleInstances) Init(logger log.Log, client Client, settings InstancesSettings) bool {

	switch asserted := settings.UnSerialize().(type) {
	case SingleInstancesSettings:
		instances.settings = asserted
	}

	return true
}
func (instances *SingleInstances) Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool {
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Single Instances")
	return true
}
func (instances *SingleInstances) Instance(id string) (instance Instance, ok bool) {
	return nil, true
}

type SingleInstance struct {

}
