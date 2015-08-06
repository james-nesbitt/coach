package main

import (
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Pull struct {
	log Log

	Nodes Nodes
	Targets []string

	Registry string

}
func (operation *Operation_Pull) Flags(flags []string) {
	operation.Registry = "registry.hub.docker.com"

}
func (operation *Operation_Pull) Run() {
	operation.log.Message("running pull operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.Targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:BUILD")
	operation.Nodes.Pull(operation.Targets, operation.Registry)
}

func (nodes *Nodes) Pull(targets []string, registry string) {
	for _, target := range nodes.GetTargets(targets) {
		target.log = nodes.log.ChildLog("NODE:"+target.Name)
		target.Pull(registry)
	}
}

func (node *Node) Pull(registry string) bool {
	if node.Do("pull") {

		image := node.GetImageName()
		tag := node.GetImageTag()

		options := docker.PullImageOptions {
			Repository: image,
			Tag: tag,
			OutputStream: os.Stdout,
			RawJSONStream:false,
		}
		if registry!="" {
			options.Registry = registry
		}

		var auth docker.AuthConfiguration
		auths, _ := docker.NewAuthConfigurationsFromDockerCfg()
		for _, auth = range auths.Configs { break }

		node.log.Message("PULLING NODE IMAGE ["+node.Name+"] : "+image)

		// ask the docker client to build the image
		err := node.client.PullImage(options, auth)

		if (err!=nil) {
			node.log.Error("NODE IMAGE FAILED ["+node.Name+"] : "+image+" => "+err.Error())
			return false
		} else {
			node.log.Message("NODE IMAGE PULLED ["+node.Name+"] : "+image)
			return true
		}

	}
	return true
}

