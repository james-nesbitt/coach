package main

import (
	"path"
	"os"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Build struct {
	log Log

	Nodes Nodes
	Targets []string

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
func (operation *Operation_Build) Run() {
	force := false
	if operation.force == true {
		force = true
	}

	operation.log.Message("running build operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.Targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:BUILD")
	operation.Nodes.Build(operation.Targets, force)
}

func (nodes *Nodes) Build(targets []string, force bool) {
	for _, target := range nodes.GetTargets(targets) {
		target.log = nodes.log.ChildLog("NODE:"+target.Name)
		target.Build(force)
	}
}

func (node *Node) Build(force bool) bool {
	if node.Do("build") && node.BuildPath!="" {

		// determine an absolute buildPath to the build, for Docker to use.
		buildPath, _ := node.conf.Path("build")
		buildPath = path.Join(buildPath, node.BuildPath)

		image := node.GetImageName()

		options := docker.BuildImageOptions{
			Name: strings.ToLower(image),
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

