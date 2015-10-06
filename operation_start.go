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
	operation.nodes.Start(operation.targets, operation.force)
}

func (nodes *Nodes) Start(targets []string, force bool) {
	for _, target := range nodes.GetTargets(targets, !force) {
		if target.node.Do("start") {
			for _, instance := range target.instances {
				instance.Start(force)
			}
		}
	}
}

func (node *Node) Start(filters []string, force bool) {
	if node.Do("start") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances(false)
		} else {
			instances = node.FilterInstances(filters, false)
		}
		for _, instance := range instances {
			if instance.active {
				instance.Start(force)
			}
		}

	}
}

func (instance *Instance) Start(force bool) bool {

	// Convert the node data into docker data (transform node keys to container IDs for things like Links & VolumesFrom)
	id := instance.GetContainerName()
	HostConfig := instance.HostConfig

	// ask the docker client to start the instance container
	err := instance.Node.client.StartContainer(id, &HostConfig)

	if err!= nil {
		instance.Node.log.Error("Failed to start node container :"+id+" => "+err.Error())
		return false
	} else {
		instance.Node.log.Message("Node instance started: "+id)
		return true
	}
}
