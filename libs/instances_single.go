package libs

import (
	"github.com/james-nesbitt/coach/log"
)

const (
	INSTANCE_SINGLE_ID          = "single"
	INSTANCE_SINGLE_MACHINENAME = ""
)

type SingleInstancesSettings struct {
	Name string
}

func (settings SingleInstancesSettings) Settings() interface{} {
	return settings
}

type SingleInstances struct {
	machineName string
	log         log.Log
	settings    SingleInstancesSettings
	instance    SingleInstance
	client      Client
}

func (instances *SingleInstances) MachineName() string {
	return instances.machineName
}
func (instances *SingleInstances) Init(logger log.Log, machineName string, client Client, settings InstancesSettings) bool {
	instances.log = logger
	instances.machineName = machineName

	instances.client = client
	switch asserted := settings.Settings().(type) {
	case SingleInstancesSettings:
		instances.settings = asserted
	}

	return true
}
func (instances *SingleInstances) Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool {
	instances.log = logger
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Single Instances")

	instances.instance = SingleInstance{}
	instances.instance.Init(logger, INSTANCE_SINGLE_ID, instances.MachineName(), client, true)

	instances.log.Debug(log.VERBOSITY_DEBUG_WOAH, "Created single instance", instances.instance)

	return true
}

func (instances *SingleInstances) Client() InstancesClient {
	return instances.client.InstancesClient(instances)
}
func (instances *SingleInstances) Instance(id string) (instance Instance, ok bool) {
	return Instance(&instances.instance), id == INSTANCE_SINGLE_ID || id == "" // return the single instance, regardless of filters
}
func (instances *SingleInstances) InstancesOrder() []string {
	return []string{INSTANCE_SINGLE_ID}
}

// Give a filterable instances for this instances object
func (instances *SingleInstances) FilterableInstances() (FilterableInstances, bool) {
	filterableInstances := BaseFilterableInstances{}
	filterableInstances.Init(Instances(instances), instances.InstancesOrder())
	return FilterableInstances(&filterableInstances), true
}

type SingleInstance struct {
	BaseInstance
}
