package libs

import (
	"github.com/james-nesbitt/coach-tools/log"
)

type Instance interface {
	Id() string
	MachineName() string

	Can(action string) bool

	Log() log.Log
	Client() InstanceClient
}

type BaseInstance struct {
	id string
	machineName string

	log log.Log
	client InstanceClient
}
func (instance *BaseInstance) Id() string {
	return instance.id
}
func (instance *BaseInstance) MachineName() string {
	return instance.machineName
}
func (instance *BaseInstance) Can(action string) bool {
	return true
}
func (instance *BaseInstance) Log() log.Log {
	return instance.log
}
func (instance *BaseInstance) Client() InstanceClient {
	return instance.client
}