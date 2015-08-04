package main

import (
	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Create struct {
	log Log

	Nodes Nodes
	Targets []string

	force bool
}
func (operation *Operation_Create) Flags(flags []string) {
	for _, flag := range flags {
		switch flag {
			case "-f":
				fallthrough
			case "--force":
				operation.force = true
		}
	}
}
func (operation *Operation_Create) Run() {
	force := false
	if operation.force == true {
		force = true
	}

	operation.log.Message("running create operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.Targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:CREATE")
	operation.Nodes.Create(operation.Targets, true, force)
}

func (nodes *Nodes) Create(targets []string, onlyActive bool, force bool) {
	for _, target := range nodes.GetTargets(targets) {
		target.Create([]string{}, onlyActive, force)
	}
}

func (node *Node) Create(filters []string, onlyActive bool, force bool) {
	var instances []*Instance

	if len(filters)==0 {
		instances = node.GetInstances(onlyActive)
	} else {
		instances = node.FilterInstances(filters, onlyActive)
	}
	for _, instance := range instances {
		instance.Create([]string{}, force)
	}
}

/**
 * Create a container for a node
 */
func (instance *Instance) Create(overrideCmd []string, force bool) bool {
	node := instance.Node

	if node.Do("create") {

		/**
			* Transform node data, into a format that can be used
			* for the actual Docker call.  This involves transforming
			* the node keys into docker container ids, for things like
			* the name, Links, VolumesFrom etc
			*/
		name := instance.GetContainerName()
		Config := instance.Config
		HostConfig := instance.HostConfig

		Config.Image = node.GetImageName()

		if len(overrideCmd)>0 {
			Config.Cmd = overrideCmd
		}

		// ask the docker client to create a container for this instance
		options := docker.CreateContainerOptions{
			Name:name,
			Config:&Config,
			HostConfig:&HostConfig,
		}

		container, err := instance.Node.client.CreateContainer( options )

		if (err!=nil) {

			/**
			 * There is a weird bug with the library, where sometimes it
			 * reports a missing image error, and yet it still creates the
			 * container.  It is not clear if this failure occurs in the
			 * remote API, or in the dockerclient library.
			 */
			if err.Error()=="no such image" {
				if container, ok := instance.GetContainer(false); ok {
					instance.Node.log.Message("CREATED INSTANCE CONTAINER ["+name+" FROM "+Config.Image+"] => "+container.ID)
					instance.Node.log.Warning("Docker created the container, but reported an error due to a 'missing image'")
					return true
				}
			}

			instance.Node.log.Error("FAILED TO CREATE INSTANCE CONTAINER ["+name+" FROM "+Config.Image+"] =>"+err.Error())
			return false
		} else {
			instance.Node.log.Message("CREATED INSTANCE CONTAINER ["+name+" FROM "+Config.Image+"] => "+container.ID)
			return true
		}

	}
	return false
}

