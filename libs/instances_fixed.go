package libs

import (
	"github.com/james-nesbitt/coach-tools/log"
)

type FixedInstancesSettings struct {
	Names []string 
}
func (settings FixedInstancesSettings) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&settings.Names)
}
func (settings FixedInstancesSettings) UnSerialize() interface{} {
	return settings
}

type FixedInstances struct {
	log log.Log
	settings FixedInstancesSettings
	instances map[string]*FixedInstance
}
func (instances *FixedInstances) Init(logger log.Log, client Client, settings InstancesSettings) bool {

	switch asserted := settings.UnSerialize().(type) {
	case FixedInstancesSettings:
		instances.settings = asserted
	}

	return true
}
func (instances *FixedInstances) Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool {
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Fixed Instances")
	return true
}
func (instances *FixedInstances) Instance(id string) (instance Instance, ok bool) {

	return nil, true
}


type FixedInstance struct {

}