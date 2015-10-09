package main

type Operation_Pause struct {
	log Log

	nodes Nodes
	targets []string

	force bool
}
func (operation *Operation_Pause) Flags(flags []string) {

}

func (operation *Operation_Pause) Help(topics []string) {
	operation.log.Note(`Operation: PAUSE

Coach will attempt to pause any target containers.

SYNTAX:
    $/> coach {targets} pause

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
  - This operation processed only nodes with the "start" access.  This excludes build, volume and command containers

NOTE:
`)
}

func (operation *Operation_Pause) Run() {
	operation.nodes.Pause(operation.targets)
}

func (nodes *Nodes) Pause(targets []string) {
	for _, target := range nodes.GetTargets(targets) {
		if target.node.Do("start") {
			for _, instance := range target.instances {
				if instance.HasContainer(true) {
					instance.Pause(false)
				}
			}
		}
	}
}

func (node *Node) Pause(filters []string, force bool) {
	if node.Do("start") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances()
		} else {
			instances = node.FilterInstances(filters)
		}

		for _, instance := range instances {
			if instance.HasContainer(true) {
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
