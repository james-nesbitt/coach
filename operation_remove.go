package main

import (
	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Remove struct {
	log Log

	Nodes Nodes
	Targets []string

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

func (operation *Operation_Remove) Help() {

}

func (operation *Operation_Remove) Run() {
	force := false
	if operation.force == true {
		force = true
	}

	operation.log.Message("running remove operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.Targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:REMOVE")
	operation.Nodes.Remove(operation.Targets, force)
}

func (nodes *Nodes) Remove(targets []string, force bool) {
	for _, target := range nodes.GetTargets(targets) {
		target.Remove([]string{}, force)
	}
}

func (node *Node) Remove(filters []string, force bool) {
	if node.Do("create") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances(true)
		} else {
			instances = node.FilterInstances(filters, true)
		}
		for _, instance := range instances {
			instance.Remove(force)
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
		instance.Node.log.Error("FAILED TO REMOVE INSTANCE CONTAINER ["+name+"] =>"+err.Error())
		return false
	} else {
		instance.Node.log.Message("REMOVED INSTANCE CONTAINER ["+name+"] ")
		return true
	}

	return false
}

