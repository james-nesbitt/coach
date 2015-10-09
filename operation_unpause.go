package main

type Operation_Unpause struct {
	log Log

	nodes Nodes
	targets []string

	force bool
}
func (operation *Operation_Unpause) Flags(flags []string) {

}

func (operation *Operation_Unpause) Help(topics []string) {
	operation.log.Note(`Operation: UNPAUSE

Coach will attempt to unpause target node containers.

SYNTAX:
    $/> coach {targets} unpause

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
  - This operation processed only nodes with the "start" access.  This excludes build, volume and command containers.
`)
}

func (operation *Operation_Unpause) Run() {
	operation.nodes.Unpause(operation.targets, operation.force)
}

func (nodes *Nodes) Unpause(targets []string, force bool) {
	for _, target := range nodes.GetTargets(targets) {
		if target.node.Do("start") {

			for _, instance := range target.instances {
				if instance.HasContainer(true) {
					instance.Unpause(force)
				}
			}

		}
	}
}

func (node *Node) Unpause(filters []string, force bool) {
	if node.Do("start") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances()
		} else {
			instances = node.FilterInstances(filters)
		}

		for _, instance := range instances {
			if instance.HasContainer(true) {
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
