package libs

import (
	"github.com/james-nesbitt/coach-tools/log"
	"math/rand"
)

const (
	INSTANCE_RANDOM_MACHINENAMESUFFIXLENGTH = 8
	INSTANCE_RANDOM_MACHINENAMECHARS        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

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
	return true
}
func (instances *TemporaryInstances) Instance(id string) (instance Instance, ok bool) {

	return nil, false
}

// Give a filterable instances for this instances object
func (instances *TemporaryInstances) FilterableInstances() (FilterableInstances, bool) {
	filterableInstances := BaseFilterableInstances{Instances: Instances(instances), filters: []string{}}
	return FilterableInstances(&filterableInstances), true
}

type TemporaryInstance struct {
	Id          string
	MachineName string
}

func (instance *TemporaryInstance) randMachineName() string {
	b := make([]byte, INSTANCE_RANDOM_MACHINENAMESUFFIXLENGTH)
	for i := range b {
		b[i] = INSTANCE_RANDOM_MACHINENAMECHARS[rand.Int63()%int64(len(INSTANCE_RANDOM_MACHINENAMECHARS))]
	}
	return instance.Id + string(b)
}
