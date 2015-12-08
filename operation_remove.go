package main

import (
	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Remove struct {
	log Log

	nodes Nodes
	targets []string

	force bool
}
func (operation *Operation_Remove) Flags(flags []string) {
	for _, flag := range flags {
		switch flag {
			case "-f":
				fallthrough
			case "--force":
				operation.force = true
		}
	}
}

func (operation *Operation_Remove) Help(topics []string) {
	operation.log.Note(`Operation: REMOVE

Coach will attempt to remove all target node containers.

SYNTAX:
	$/> coach {targets} remove

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
	- only nodes with the "create" access are processed.  This excludes build and command nodes
`)
}

func (operation *Operation_Remove) Run() {
	force := false
	if operation.force == true {
		force = true
	}

	operation.log.Info("running remove operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:REMOVE")
	operation.nodes.Remove(operation.targets, force)
}

func (nodes *Nodes) Remove(targets []string, force bool) {
	for _, target := range nodes.GetTargets(targets) {
		if target.node.Do("create") {
			for _, instance := range target.instances {
				if instance.HasContainer(false) {
					instance.Remove(force)
				} else {
					instance.Node.log.Info(instance.Node.Name+": Skipping remove as this instance ["+instance.Name+"] has no container ")
				}
			}
		}
	}
}

func (node *Node) Remove(filters []string, force bool) {
	if node.Do("create") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances()
		} else {
			instances = node.FilterInstances(filters)
		}

		for _, instance := range instances {
			if instance.HasContainer(false) {
				instance.Remove(force)
			} else {
				node.log.Info(node.Name+": Skipping remove as this instance ["+instance.Name+"] has no container ")
			}
		}

	}
}

func (instance *Instance) Remove(force bool) bool {

	name := instance.GetContainerName()
	options := docker.RemoveContainerOptions{
		ID: name,
	}

	// ask the docker client to remove the instance container
	err := instance.Node.client.RemoveContainer(options)

	if (err!=nil) {
		instance.Node.log.Error(instance.Node.Name+": Failed to remove instance container ["+name+"] =>"+err.Error())
		return false
	} else {
		instance.Node.log.Message(instance.Node.Name+": Removed instance container ["+name+"] ")
		return true
	}

	return false
}

