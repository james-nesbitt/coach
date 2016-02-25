package libs

import (
	"github.com/james-nesbitt/coach/log"
)

const (
	INSTANCES_NULL_MACHINENAME = "$null"
)

type NullInstancesSettings struct{}

func (settings NullInstancesSettings) Settings() interface{} {
	return settings
}

type NullInstances struct {
	log log.Log
}

func (instances *NullInstances) Init(logger log.Log, machineName string, client Client, settings InstancesSettings) bool {
	instances.log = logger
	return true
}
func (instances *NullInstances) Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool {
	instances.log = logger
	return true
}

func (instances *NullInstances) Client() InstancesClient {
	instances.log.Info("Tried to retrieve a client for a Null Instances.")
	return nil
}
func (instances *NullInstances) MachineName() string {
	return INSTANCES_NULL_MACHINENAME
}

// This instances type can never actually provide an instance
func (instances *NullInstances) Instance(name string) (Instance, bool) {
	instances.log.Info("Tried to retrieve a Null Instances.")
	return nil, false
}
func (instances *NullInstances) InstancesOrder() []string {
	return []string{}
}
func (instances *NullInstances) FilterableInstances() (FilterableInstances, bool) {
	return nil, false
}

// Status strings
func (instances *NullInstances) Status(logger log.Log) []string {
	return []string{}
}