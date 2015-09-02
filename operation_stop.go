package main

type Operation_Stop struct {
	log Log

	Nodes Nodes
	Targets []string

	force bool
	timeout uint
}
func (operation *Operation_Stop) Flags(flags []string) {

}

func (operation *Operation_Stop) Help(topics []string) {
	operation.log.Note(`Operation: STOP

Coach will attempt to stop target node containers.
`)
}

func (operation *Operation_Stop) Run() {
	operation.Nodes.Stop(operation.Targets, operation.force, operation.timeout)
}

func (nodes *Nodes) Stop(targets []string, force bool, timeout uint) {
	for _, target := range nodes.GetTargets(targets) {
		target.Stop([]string{}, force, timeout)
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
