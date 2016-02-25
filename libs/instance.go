package libs

import (
	"github.com/james-nesbitt/coach/log"
)

type Instance interface {
	Init(logger log.Log, id string, machineName string, client Client, isDefault bool) bool

	Id() string
	MachineName() string

	Can(action string) bool
	IsDefault() bool
	IsRunning() bool
	IsReady() bool

	Log() log.Log
	Client() InstanceClient

	Status(logger log.Log) []string
}

type BaseInstance struct {
	id          string
	machineName string

	log    log.Log
	client Client

	isDefault bool
}

func (instance *BaseInstance) Init(logger log.Log, id string, machineName string, client Client, isDefault bool) bool {
	instance.id = id
	instance.machineName = machineName
	instance.log = logger
	instance.client = client
	instance.isDefault = isDefault

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
func (instance *BaseInstance) IsDefault() bool {
	return instance.isDefault
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

func (instance *BaseInstance) Status(logger log.Log) []string {
	status := []string{}

	instanceClient := instance.Client()
	if instanceClient.IsRunning() {
		status = append(status, "RUNNING")
	} else if instanceClient.HasContainer() {
		status = append(status, "READY")
	} else {
		status = append(status, "NOT-READY")
	}

	return status
}