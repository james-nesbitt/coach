package main

import (
	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Destroy struct {
	log Log

	Nodes Nodes
	Targets []string

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
func (operation *Operation_Destroy) Run() {
	force := false
	if operation.force == true {
		force = true
	}

	operation.log.Message("running destroy operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.Targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:BUILD")
	operation.Nodes.Destroy(operation.Targets, force)
}

func (nodes *Nodes) Destroy(targets []string, force bool) {
	for _, target := range nodes.GetTargets(targets) {
		target.log = nodes.log.ChildLog("NODE:"+target.Name)
		target.Destroy(force)
	}
}

func (node *Node) Destroy(force bool) bool {
	if node.Do("build") {

		// Get the image name
		image := node.GetImageName()
		if tag := node.GetImageTag(); tag!="" && tag!="latest" {
			image +=":"+tag
		}

		options := docker.RemoveImageOptions{
			Force: force,
		}

		// ask the docker client to remove the image
		err := node.client.RemoveImageExtended(image, options)

		if (err!=nil) {
			node.log.Error("NODE DESTROY FAILED ["+node.Name+"] : "+image+" => "+err.Error())
			return false
		} else {
			node.log.Message("NODE DESTROYED ["+node.Name+"] : "+image)
			return true
		}

	}
	return true
}
