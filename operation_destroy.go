package main

import (
	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Destroy struct {
	log Log

	nodes Nodes
	targets []string

	force bool
}
func (operation *Operation_Destroy) Flags(flags []string) {
	for _, flag := range flags {
		switch flag {
			case "-f":
				fallthrough
			case "--force":
				operation.force = true
		}
	}
}

func (operation *Operation_Destroy) Help(topics []string) {
	operation.log.Note(`Operation: DESTROY

Coach will attempt to remove any built images for target nodes.  A node is
considered "buildable" if it has a Build: definition.

SYNTAX:
    $/> coach {targets} destroy

  {targets} what target nodes the operation should process ($/> coach help targets)

ACCESS:
  - this opertion will only process nodes with "build" access.  This includes only nodes with the Build: settings declared.

NOTE:
- Coach will try not to remove an image for a node that does not build, which is the common case for nodes with an image, but no build setting.  This prevents deleting shared images that are not build targets.
`)
}

func (operation *Operation_Destroy) Run() {
	force := false
	if operation.force == true {
		force = true
	}

	operation.log.Info("running destroy operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:BUILD")
	operation.nodes.Destroy(operation.targets, force)
}

func (nodes *Nodes) Destroy(targets []string, force bool) {
	for _, target := range nodes.GetTargets(targets) {
		target.node.log = nodes.log.ChildLog("NODE:"+target.node.Name)
		target.node.Destroy(force)
	}
}

func (node *Node) Destroy(force bool) bool {
	if node.Do("build") {

		// Get the image name
		image, tag := node.GetImageName()
		if tag!="" {
			image +=":"+tag
		}

		options := docker.RemoveImageOptions{
			Force: force,
		}

		// ask the docker client to remove the image
		err := node.client.RemoveImageExtended(image, options)

		if (err!=nil) {
			node.log.Error(node.Name+": Node image removal failed ["+image+"] => "+err.Error())
			return false
		} else {
			node.log.Message(node.Name+": Node image was removed ["+image+"]")
			return true
		}

	} else {
		node.log.Info(node.Name+": Node has no built image, so it will not be destroyed.")
	}
	return true
}
