package libs

import (
	"math/rand"
	"github.com/james-nesbitt/coach-tools/log"
)

const (
	INSTANCE_RANDOM_MACHINENAMESUFFIXLENGTH = 8
	INSTANCE_RANDOM_MACHINENAMECHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type TemporaryInstancesSettings struct {
	Name string 
}
func (settings TemporaryInstancesSettings) UnSerialize() interface{} {
	return settings
}


type TemporaryInstances struct {
	log log.Log
	settings TemporaryInstancesSettings	
	client Client
}
func (instances *TemporaryInstances) Init(logger log.Log, client Client, settings InstancesSettings) bool {

	switch asserted := settings.UnSerialize().(type) {
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
type TemporaryInstance struct {
	Id string
	MachineName string
}
func (instance *TemporaryInstance) randMachineName() string {
  b := make([]byte, INSTANCE_RANDOM_MACHINENAMESUFFIXLENGTH)
  for i := range b {
      b[i] = INSTANCE_RANDOM_MACHINENAMECHARS[rand.Int63() % int64(len(INSTANCE_RANDOM_MACHINENAMECHARS))]
  }
  return instance.Id+string(b)
}
