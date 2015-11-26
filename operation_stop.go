package main

type Operation_Stop struct {
	log Log

	nodes Nodes
	targets []string

	force bool
	timeout uint
}
func (operation *Operation_Stop) Flags(flags []string) {
  operation.timeout=10

	remainingFlags := []string{}

  flagLoop:
		for index:=0; index<len(flags); index++ {
			flag:= flags[index]

			switch flag {
				case "-q":
					fallthrough
				case "--quick":
					operation.timeout=1

			}
		}


}

func (operation *Operation_Stop) Help(topics []string) {
	operation.log.Note(`Operation: STOP

Coach will attempt to stop target node containers.

SYNTAX:
    $/> coach {targets} stop

    $/> coach {targets} stop --quick
      - makes docker stop the containers with --time=1

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
  - This operation processed only nodes with the "start" access.  This excludes build, volume and command containers.

`)
}

func (operation *Operation_Stop) Run() {
	operation.log.Info("running stop operation")	
	operation.nodes.Stop(operation.targets, operation.force, operation.timeout)
}

func (nodes *Nodes) Stop(targets []string, force bool, timeout uint) {
	for _, target := range nodes.GetTargets(targets) {
		if target.node.Do("start") {
			for _, instance := range target.instances {
				if instance.HasContainer(true) {
					instance.Stop(force, timeout)
				} else {
					target.node.log.Info(target.node.Name+": instance has no running container to stop ["+instance.Name+"]")
				}
			}
		} else {
			target.node.log.Info(target.node.Name+": is not a stoppable node")
		}
	}
}

func (node *Node) Stop(filters []string, force bool, timeout uint) {
	if node.Do("start") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances()
		} else {
			instances = node.FilterInstances(filters)
		}

		for _, instance := range instances {
			if instance.HasContainer(true) {
				instance.Stop(force, timeout)
			} else {
				node.log.Message(node.Name+": instance has no running container to stop ["+instance.Name+"]")
			}

		}

	} else {
		node.log.Info(node.Name+": is not a stoppable node")
	}
}


func (instance *Instance) Stop(force bool, timeout uint) bool {

	id := instance.GetContainerName()
	err := instance.Node.client.StopContainer(id, timeout)

	if err!= nil {
		instance.Node.log.Error(instance.Node.Name+": Failed to stop node container ["+id+"] => "+err.Error())
		return false
	} else {
		instance.Node.log.Message(instance.Node.Name+": Node instance stopped ["+id+"]")
		return true
	}
}
