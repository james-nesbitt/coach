package main

import (
	"path"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Build struct {
	log Log

	nodes Nodes
	targets []string

	force bool
}
func (operation *Operation_Build) Flags(flags []string) {
	for _, flag := range flags {
		switch flag {
			case "-f":
				fallthrough
			case "--force":
				operation.force = true
		}
	}
}


func (operation *Operation_Build) Help(topics []string) {
	operation.log.Note(`Operation: BUILD

Coach will attempt to build a new docker image, for each target node that has a build setting.

The operation will look for a Build: setting inside the node, and try to find a matching Dockerfile at the suggested path.  The path can be relative to the project root, or absolute.

SYNTAX:
    $/> coach {targets} build

  {targets} what target nodes the operation should process ($/> coach help targets)

ACCESS:
  - this operation processes only nodes with the "build" access.  This includes only nodes with a Build: setting.

NOTES:
- a node that can be built should have a Build: setting, which points to a project path that contains the Dockerfile.
- while {targets} globally can specify particular node instances, that information is ignored for this operation, as images are built for all instances of a node.
`)
}

func (operation *Operation_Build) Run() {
	force := false
	if operation.force == true {
		force = true
	}

	operation.log.Message("running build operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:BUILD")
	operation.nodes.Build(operation.targets, force)
}

func (nodes *Nodes) Build(targets []string, force bool) {
	for _, target := range nodes.GetTargets(targets) {
		target.node.log = nodes.log.ChildLog("NODE:"+target.node.Name)

		// ignore target instances, and just build the node
		target.node.Build(force)
	}
}

func (node *Node) Build(force bool) bool {
	if node.Do("build") && node.BuildPath!="" {

		// determine an absolute buildPath to the build, for Docker to use.
		buildPath, _ := node.conf.Path("build")
		buildPath = path.Join(buildPath, node.BuildPath)

		image := node.GetImageName()
		if tag := node.GetImageTag(); tag!="" {
			image +=":"+tag
		}

		options := docker.BuildImageOptions{
			Name: image,
			ContextDir: buildPath,
			RmTmpContainer: true,
			OutputStream: os.Stdout,
		}

		node.log.Message("BUILDING NODE ["+node.Name+"] : "+image+" FROM "+buildPath)

		// ask the docker client to build the image
		err := node.client.BuildImage( options )

		if (err!=nil) {
			node.log.Error("NODE BUILD FAILED ["+node.Name+"] : "+buildPath+" => "+err.Error())
			return false
		} else {
			node.log.Message("NODE BUILT ["+node.Name+"] : "+image+" FROM "+buildPath)
			return true
		}

	}
	return true
}

