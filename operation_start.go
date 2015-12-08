package main

type Operation_Start struct {
	log Log

	nodes Nodes
	targets []string

	force bool
}
func (operation *Operation_Start) Flags(flags []string) {

}

func (operation *Operation_Start) Help(topics []string) {
	operation.log.Note(`Operation: START

Coach will attempt to start target node containers.

SYNTAX:
	$/> coach {targets} start

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
	- This operation processed only nodes with the "start" access.  This excludes build, volume and command containers.
`)
}

func (operation *Operation_Start) Run() {
	operation.log.Info("running start operation")
	operation.nodes.Start(operation.targets, operation.force)
}

func (nodes *Nodes) Start(targets []string, force bool) {
	for _, target := range nodes.GetTargets(targets) {
		if target.node.Do("start") {
			for _, instance := range target.instances {
				if instance.isRunning() {
					target.node.log.Info(target.node.Name+"["+instance.Name+"]: is already running.")
				} else if instance.HasContainer(false) {
					instance.Start(force)
				} else {
					target.node.log.Info(target.node.Name+"["+instance.Name+"]: has no container to start.")
				}
			}
		} else {
			target.node.log.Info(target.node.Name+": is not a startable node.")
		}
	}
}

func (node *Node) Start(filters []string, force bool) {
	if node.Do("start") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances()
		} else {
			instances = node.FilterInstances(filters)
		}
		for _, instance := range instances {
			if instance.isRunning() {
				node.log.Info(node.Name+": instance is already running ["+instance.Name+"]")
			} else if instance.HasContainer(false) {
				instance.Start(force)
			} else {
				node.log.Info(node.Name+": has no container to start ["+instance.Name+"]")
			}
		}

	} else {
		node.log.Info(node.Name+": is not a startable node.")
	}
}

func (instance *Instance) Start(force bool) bool {

	// Convert the node data into docker data (transform node keys to container IDs for things like Links & VolumesFrom)
	id := instance.GetContainerName()
	HostConfig := instance.HostConfig

	// ask the docker client to start the instance container
	err := instance.Node.client.StartContainer(id, &HostConfig)

	if err!= nil {
		instance.Node.log.Error(instance.Node.Name+": Failed to start node container ["+id+"] => "+err.Error())
		return false
	} else {
		instance.Node.log.Message(instance.Node.Name+": Node instance started ["+id+"]")
		return true
	}
}
