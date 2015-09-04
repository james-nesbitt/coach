package main

type Operation_Stop struct {
	log Log

	nodes Nodes
	targets []string

	force bool
	timeout uint
}
func (operation *Operation_Stop) Flags(flags []string) {

}

func (operation *Operation_Stop) Help(topics []string) {
	operation.log.Note(`Operation: STOP

Coach will attempt to stop target node containers.

SYNTAX:
    $/> coach {targets} stop

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
  - This operation processed only nodes with the "start" access.  This excludes build, volume and command containers.

`)
}

func (operation *Operation_Stop) Run() {
	operation.nodes.Stop(operation.targets, operation.force, operation.timeout)
}

func (nodes *Nodes) Stop(targets []string, force bool, timeout uint) {
	for _, target := range nodes.GetTargets(targets, !force) {
		if target.node.Do("start") {
			for _, instance := range target.instances {
				instance.Stop(force, timeout)
			}
		}
	}
}

func (node *Node) Stop(filters []string, force bool, timeout uint) {
	if node.Do("start") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances(true)
		} else {
			instances = node.FilterInstances(filters, true)
		}
		for _, instance := range instances {
			if instance.active {
				instance.Stop(force, timeout)
			}
		}

	}
}


func (instance *Instance) Stop(force bool, timeout uint) bool {

	id := instance.GetContainerName()
	err := instance.Node.client.StopContainer(id, timeout)

	if err!= nil {
		instance.Node.log.Error("Failed to stop node container :"+id+" => "+err.Error())
		return false
	} else {
		instance.Node.log.Message("Node instance stopped: "+id)
		return true
	}
}
