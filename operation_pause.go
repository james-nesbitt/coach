package main

type Operation_Pause struct {
	log Log

	Nodes Nodes
	Targets []string

	force bool
}
func (operation *Operation_Pause) Flags(flags []string) {

}

func (operation *Operation_Pause) Help() {

}

func (operation *Operation_Pause) Run() {
	operation.Nodes.Pause(operation.Targets)
}

func (nodes *Nodes) Pause(targets []string) {
	for _, target := range nodes.GetTargets(targets) {
		target.Pause([]string{}, false)
	}
}

func (node *Node) Pause(filters []string, force bool) {
	if node.Do("start") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances(false)
		} else {
			instances = node.FilterInstances(filters, false)
		}
		for _, instance := range instances {
			if instance.active {
				instance.Pause(force)
			}
		}

	}
}

func (instance *Instance) Pause(force bool) bool {

	id := instance.GetContainerName()

	err := instance.Node.client.PauseContainer( id )
	if err!=nil {
		instance.Node.log.Error("FAILED TO PAUSE INSTANCE CONTAINER ["+id+"] =>"+err.Error())
		return false
	} else {
		instance.Node.log.Message("PAUSED INSTANCE CONTAINER ["+id+"]")
		return true
	}

}
