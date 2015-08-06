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
	operation.Registry = "https://index.docker.io/v1/"
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

		var auth docker.AuthConfiguration
		auths, _ := docker.NewAuthConfigurationsFromDockerCfg()
		if auths==nil {
			node.log.Warning("You have no local login credentials for any repo")
			auth = docker.AuthConfiguration{}
			#options.Registry = "https://index.docker.io/v1/"
		} else {
			if registry!="" {
				auth, _ = auths.Configs[registry]
				options.Registry = registry
			}
			if auth.Username=="" {
				for registry, regauth := range auths.Configs {
					options.Registry = registry
					auth = regauth
					break
				}
			}
		}

		node.log.Message("PULLING NODE IMAGE ["+node.Name+"] FROM SERVER ["+options.Registry+"] USING AUTH ["+auth.Username+"] : "+image)
		node.log.DebugObject( LOG_SEVERITY_DEBUG_LOTS, "AUTH USED: ", map[string]string{"Username":auth.Username, "Email":auth.Email, "ServerAdddress":auth.ServerAddress})

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

