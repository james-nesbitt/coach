package libs

import (
	"github.com/james-nesbitt/coach-tools/log"
	"strconv"
)

const (
	INSTANCES_SCALED_DEFAULT_INITIAL = 3
	INSTANCES_SCALED_DEFAULT_MAXIMUM = 6
)

type ScaledInstancesSettings struct {
	Initial int `json:"Initial,omitempty" yaml:"Initial,omitempty"`
	Maximum int `json:"Maximum,omitempty" yaml:"Maximum,omitempty"`
}

func (settings ScaledInstancesSettings) Settings() interface{} {
	return settings
}

type ScaledInstances struct {
	BaseInstances

	settings ScaledInstancesSettings
}

func (instances *ScaledInstances) Init(logger log.Log, machineName string, client Client, settings InstancesSettings) bool {
	instances.BaseInstances.Init(logger, machineName, client, settings)

	switch asserted := settings.Settings().(type) {
	case ScaledInstancesSettings:
		// sanity check
		if asserted.Initial < 0 {
			asserted.Maximum = INSTANCES_SCALED_DEFAULT_INITIAL
			logger.Warning("Scaled instances Initial count was invalid.  Switching to the default value")
		}
		if asserted.Maximum < 1 {
			asserted.Maximum = INSTANCES_SCALED_DEFAULT_MAXIMUM
			logger.Warning("Scaled instances Maximum was invalid.  Switching to the default value")
		}
		if asserted.Initial > asserted.Maximum {
			asserted.Maximum = asserted.Initial
			logger.Warning("Scaled instances Maximum is lower than the Initial value.  Using Initial as Maximum.")
		}

		instances.settings = asserted
	default:
		instances.settings = ScaledInstancesSettings{Initial: INSTANCES_SCALED_DEFAULT_INITIAL, Maximum: INSTANCES_SCALED_DEFAULT_MAXIMUM}
	}

	return true
}
func (instances *ScaledInstances) Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool {
	instances.log = logger
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Scaled Instances")

	for i := 0; i <= instances.settings.Maximum; i++ {
		name := strconv.Itoa(int(i))
		machineName := instances.MakeId(i)
		instance := Instance(&ScaledInstance{})

		if instance.Init(logger.MakeChild(name), name, machineName, client) {
			instances.instancesMap[name] = instance
			instances.instancesOrder = append(instances.instancesOrder, name)
		}
	}

	return true
}

func (instances *ScaledInstances) MakeId(index int) string {
	return instances.MachineName() + "_" + strconv.Itoa(int(index))
}

// Give a filterable instances for this instances object
func (instances *ScaledInstances) FilterableInstances() (FilterableInstances, bool) {
	filterableInstances := BaseFilterableInstances{Instances: Instances(instances), filters: []string{}}
	return FilterableInstances(&filterableInstances), true
}

type ScaledInstance struct {
	BaseInstance
}
