package main

type Operation_Unpause struct {
	log Log

	Nodes Nodes
	Targets []string

	force bool
}
func (operation *Operation_Unpause) Flags(flags []string) {

}
func (operation *Operation_Unpause) Run() {
	operation.Nodes.Unpause(operation.Targets)
}

func (nodes *Nodes) Unpause(targets []string) {
	for _, target := range nodes.GetTargets(targets) {
		target.Unpause([]string{}, false)
	}
}

func (node *Node) Unpause(filters []string, force bool) {
	if node.Do("start") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances(false)
		} else {
			instances = node.FilterInstances(filters, false)
		}
		for _, instance := range instances {
			if instance.active {
				instance.Unpause(force)
			}
		}

	}
}

func (instance *Instance) Unpause(force bool) bool {
	id := instance.GetContainerName()

	err := instance.Node.client.UnpauseContainer( id )
	if err!=nil {
		instance.Node.log.Error("FAILED TO UNPAUSE INSTANCE CONTAINER ["+id+"] =>"+err.Error())
		return false
	} else {
		instance.Node.log.Message("UNPAUSED INSTANCE CONTAINER ["+id+"]")
		return true
	}
}
