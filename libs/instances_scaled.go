package libs

import (
	"strconv"
	"github.com/james-nesbitt/coach-tools/log"
)

type ScaledInstancesSettings struct {
	Initial int  `json:"Initial,omitempty" yaml:"Initial,omitempty"`
	Maximum int  `json:"Maximum,omitempty" yaml:"Maximum,omitempty"` 
}
func (settings ScaledInstancesSettings) UnSerialize() interface{} {
	return settings
}

type ScaledInstances struct {
	BaseInstances

	settings ScaledInstancesSettings
}
func (instances *ScaledInstances) Init(logger log.Log, client Client, settings InstancesSettings) bool {

	switch asserted := settings.UnSerialize().(type) {
	case ScaledInstancesSettings:
		instances.settings = asserted
	default:
		instances.settings = ScaledInstancesSettings{Initial:3, Maximum: 3}
	}
	return true
}
func (instances *ScaledInstances) Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool {
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Scaled Instances")
	return true
}
func (instances *ScaledInstances) Instance(id string) (instance Instance, ok bool) {

	return nil, true
}
func (instances *ScaledInstances) MakeId(index int) string {
	return strconv.Itoa(index)
}