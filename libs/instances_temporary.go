package libs

import (
	"github.com/james-nesbitt/coach/log"

	"math/rand"
	"time"
)

const (
	INSTANCE_RANDOM_MACHINENAMESUFFIXLENGTH = 8
	INSTANCE_RANDOM_MACHINENAMECHARS        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type TemporaryInstancesSettings struct {
	Name string
}

func (settings TemporaryInstancesSettings) Settings() interface{} {
	return settings
}

type TemporaryInstances struct {
	BaseInstances

	settings TemporaryInstancesSettings
}

func (instances *TemporaryInstances) Init(logger log.Log, machineName string, client Client, settings InstancesSettings) bool {
	instances.BaseInstances.Init(logger, machineName, client, settings)

	switch asserted := settings.Settings().(type) {
	case TemporaryInstancesSettings:
		instances.settings = asserted
	}

	return true
}
func (instances *TemporaryInstances) Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool {
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Temporary Instances")

	// look for existing instances, which can come from persistant runs, or run breaks
	for _, instanceId := range instances.client.InstancesClient(instances).InstancesFound(logger) {
		instances.Instance(instanceId)
	}

	return true
}
func (instances *TemporaryInstances) Instance(id string) (instance Instance, ok bool) {
	if instance, ok = instances.instancesMap[id]; ok {
		return instance, ok
	}

	if id == "" {
		id = instances.randMachineName()
	}

	machineName := instances.MachineName() + "_" + id

	instance = Instance(&TemporaryInstance{})
	instance.Init(instances.log.MakeChild(id), id, machineName, instances.client, false)

	instances.instancesMap[id] = instance
	instances.instancesOrder = append(instances.instancesOrder, id)
	return instance, true
}

// Status strings
func (instances *TemporaryInstances) Status(logger log.Log) []string {
	status := []string{}

	for _, id := range instances.InstancesOrder() {
		instance, _ := instances.Instance(id)
		instanceClient := instance.Client()
		if instanceClient.IsRunning() {
			status = append(status, "ACTIVE:"+instance.MachineName())
		} else if instanceClient.HasContainer() {
			status = append(status, "ORPHAN:"+instance.MachineName())
		}
	}

	return status
}

// Give a filterable instances for this instances object
func (instances *TemporaryInstances) FilterableInstances() (FilterableInstances, bool) {
	filterableInstances := BaseFilterableInstances{}
	filterableInstances.Init(Instances(instances), []string{""}) // default to dealing with a new temporary instance
	return FilterableInstances(&filterableInstances), true
}

func (instances *TemporaryInstances) randMachineName() string {
	b := make([]byte, INSTANCE_RANDOM_MACHINENAMESUFFIXLENGTH)
	for i := range b {
		b[i] = INSTANCE_RANDOM_MACHINENAMECHARS[rand.Int63()%int64(len(INSTANCE_RANDOM_MACHINENAMECHARS))]
	}
	return string(b)
}

type TemporaryInstance struct {
	BaseInstance
}
