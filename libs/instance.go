package libs

import (
	"github.com/james-nesbitt/coach/log"
)

type Instance interface {
	Init(logger log.Log, id string, machineName string, client Client) bool

	Id() string
	MachineName() string

	Can(action string) bool
	IsRunning() bool
	IsReady() bool

	Log() log.Log
	Client() InstanceClient
}

type BaseInstance struct {
	id          string
	machineName string

	log    log.Log
	client Client
}

func (instance *BaseInstance) Init(logger log.Log, id string, machineName string, client Client) bool {
	instance.id = id
	instance.machineName = machineName
	instance.log = logger
	instance.client = client

	return true
}
func (instance *BaseInstance) Prepare() bool {
	return true
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
func (instance *BaseInstance) IsReady() bool {
	return instance.Client().HasContainer()
}
func (instance *BaseInstance) IsRunning() bool {
	return instance.Client().IsRunning()
}

func (instance *BaseInstance) Log() log.Log {
	return instance.log
}
func (instance *BaseInstance) Client() InstanceClient {
	return instance.client.InstanceClient(instance)
}
