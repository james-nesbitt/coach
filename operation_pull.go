package main

import (
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Pull struct {
	log Log

	Nodes Nodes
	Targets []string

	Repository string
	Registry string

}
func (operation *Operation_Pull) Flags(flags []string) {
	operation.Repository = "registry.hub.docker.com"

}
func (operation *Operation_Pull) Run() {
	operation.log.Message("running pull operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.Targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:BUILD")
	operation.Nodes.Pull(operation.Targets, operation.Repository, operation.Registry)
}

func (nodes *Nodes) Pull(targets []string, repository string, registry string) {
	for _, target := range nodes.GetTargets(targets) {
		target.log = nodes.log.ChildLog("NODE:"+target.Name)
		target.Pull(repository, registry)
	}
}

func (node *Node) Pull(repository string, registry string) bool {
	if node.Do("pull") {

		image := node.Config.Image

		options := docker.PullImageOptions {
			Tag: image,
			OutputStream: os.Stdout,
			RawJSONStream:false,
		}
		if repository!="" {
			options.Repository = repository
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

