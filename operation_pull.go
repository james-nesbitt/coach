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
		if tag := node.GetImageTag(); tag!="" && tag!="latest" {
			image +=":"+tag
		}

		options := docker.PullImageOptions {
			Repository: image,
			OutputStream: os.Stdout,
			RawJSONStream:false,
		}

		tag:=node.GetImageTag();
		if tag!="" {
			options.Tag = tag
		}

		var auth docker.AuthConfiguration
// 		var ok bool
		//options.Registry = "https://index.docker.io/v1/"

// 		auths, _ := docker.NewAuthConfigurationsFromDockerCfg()
// 		if auth, ok = auths.Configs[registry]; ok {
// 			options.Registry = registry
// 		} else {
// 			node.log.Warning("You have no local login credentials for any repo. Defaulting to no login.")
			auth = docker.AuthConfiguration{}
			options.Registry = "https://index.docker.io/v1/"
// 		}

		node.log.Message("PULLING NODE IMAGE ["+node.Name+"] FROM SERVER ["+options.Registry+"] USING AUTH ["+auth.Username+"] : "+image+":"+tag)
		node.log.DebugObject( LOG_SEVERITY_DEBUG_LOTS, "AUTH USED: ", map[string]string{"Username":auth.Username, "Password":auth.Password, "Email":auth.Email, "ServerAdddress":auth.ServerAddress})

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
