package main

import (
	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Create struct {
	log Log

	nodes Nodes
	targets []string

	cmd []string

	nonDefault bool
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

func (operation *Operation_Create) Help(topics []string) {
	operation.log.Note(`Operation: CREATE

Coach will attempt to create any node containers that should be active.

Syntax:
    $/> coach {targets} create

  {targets} what target node instances the operation should process ($/> coach help targets)

Access:
  - only nodes with the "create" access are processed.  This excludes build and command nodes
`)
}

func (operation *Operation_Create) Run() {
	operation.log.Info("running create operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:CREATE")
	operation.nodes.Create(operation.targets, operation.cmd, !operation.nonDefault, operation.force)
}

func (nodes *Nodes) Create(targets []string, cmdOverride []string, onlyDefault bool, force bool) {
	for _, target := range nodes.GetTargets(targets) {

		if target.node.Do("create") {
			for _, instance := range target.instances {
				if instance.HasContainer(false) {
					nodes.log.Info(target.node.Name+"["+instance.Name+"]: Skipping node instance, which already has a container")
					continue
				}
				if onlyDefault && !instance.isDefault() {
					nodes.log.Info(target.node.Name+"["+instance.Name+"]: Skipping node instance, which is not created by default default")
					continue
				}

				nodes.log.Info("Creating node instance :"+target.node.Name+":"+instance.Name)
				instance.Create(cmdOverride, force)
			}
		} else {
			nodes.log.Info("Skipping 'uncreateable' node :"+target.node.Name)
		}
	}
}

func (node *Node) Create(filters []string, cmdOverride []string, onlyDefault bool, force bool) {
	if node.Do("create") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances()
		} else {
			instances = node.FilterInstances(filters)
		}

		for _, instance := range instances {
			if instance.HasContainer(false) {
				node.log.Info(node.Name+"["+instance.Name+"]: Skipping node instance, which already has a container")
				continue
			}
			if onlyDefault && !instance.isDefault() {
				node.log.Info(node.Name+"["+instance.Name+"]: Skipping node instance, which is not created by default default")	
				continue
			}
			instance.Create(cmdOverride, force)
		}

	}
}

/**
 * Create a container for a node
 */
func (instance *Instance) Create(overrideCmd []string, force bool) bool {
	node := instance.Node

	/**
		* Transform node data, into a format that can be used
		* for the actual Docker call.  This involves transforming
		* the node keys into docker container ids, for things like
		* the name, Links, VolumesFrom etc
		*/
	name := instance.GetContainerName()
	Config := instance.Config
	HostConfig := instance.HostConfig

	image, tag := node.GetImageName()
	if tag!="" && tag!="latest" {
		image +=":"+tag
	}
	Config.Image = image

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
	cache.refresh(false, true)

	if (err!=nil) {

		instance.Node.log.DebugObject(LOG_SEVERITY_DEBUG, "CREATE FAIL CONTAINERS: ", err)

		/**
			* There is a weird bug with the library, where sometimes it
			* reports a missing image error, and yet it still creates the
			* container.  It is not clear if this failure occurs in the
			* remote API, or in the dockerclient library.
			*/

		if err.Error()=="no such image" && instance.HasContainer(false) {
			instance.Node.log.Message(instance.Node.Name+": Created instance container ["+name+" FROM "+Config.Image+"] => "+container.ID)
			instance.Node.log.Warning("Docker created the container, but reported an error due to a 'missing image'.  This is a known bug, that can be ignored")
			return true
		}

		instance.Node.log.Error(instance.Node.Name+": Failed to create instance container ["+name+" FROM "+Config.Image+"] => "+err.Error())
		return false
	} else {
		instance.Node.log.Message(instance.Node.Name+": Created instance container ["+name+"] => "+container.ID[:12])
		return true
	}
}
