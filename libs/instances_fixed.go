package libs

import (
	"github.com/james-nesbitt/coach/log"
)

type FixedInstancesSettings struct {
	Names []string
}

func (settings FixedInstancesSettings) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&settings.Names)
}
func (settings FixedInstancesSettings) Settings() interface{} {
	return settings
}

type FixedInstances struct {
	BaseInstances

	log      log.Log
	settings FixedInstancesSettings
}

func (instances *FixedInstances) Init(logger log.Log, machineName string, client Client, settings InstancesSettings) bool {
	instances.BaseInstances.Init(logger, machineName, client, settings)

	switch asserted := settings.Settings().(type) {
	case FixedInstancesSettings:
		instances.settings = asserted
	}

	return true
}
func (instances *FixedInstances) Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool {
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Fixed Instances")

	for _, name := range instances.settings.Names {
		machineName := instances.MachineName() + "_" + name

		instance := Instance(&FixedInstance{})
		if instance.Init(logger.MakeChild(name), name, machineName, client, true) {
			instances.instancesMap[name] = instance
			instances.instancesOrder = append(instances.instancesOrder, name)
		}
	}

	return true
}

// Give a filterable instances for this instances object
func (instances *FixedInstances) FilterableInstances() (FilterableInstances, bool) {
	filterableInstances := BaseFilterableInstances{}
	filterableInstances.Init(Instances(instances), instances.InstancesOrder())
	return FilterableInstances(&filterableInstances), true
}

type FixedInstance struct {
	BaseInstance
}
