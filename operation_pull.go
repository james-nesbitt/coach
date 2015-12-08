package main

import (
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Pull struct {
	log Log

	nodes Nodes
	targets []string

	Registry string
	force bool
}
func (operation *Operation_Pull) Flags(flags []string) {
	operation.Registry = "https://index.docker.io/v1/"
	operation.force = false

	for index:=0; index<len(flags); index++ {
		flag:= flags[index]

		switch flag {
			case "-f":
				fallthrough
			case "--force":
				operation.force = true
			case "-r":
			  fallthrough
			case "--registry":
				index++
				operation.Registry = flags[index]
		}
	}	
}

func (operation *Operation_Pull) Help(topics []string) {
	operation.log.Note(`Operation: PULL

Coach will attempt to pull any node images, for nodes that have no build settings.


SYNTAX:
    $/> coach {targets} pull

  {targets} what target nodes the operation should process ($/> coach help targets)

ACCESS:
  - this operation processes only nodes with the "pull" access.  This includes only nodes without a Build: setting, but a Config: Image: setting.

NOTES:
  - Nodes that have build settings will not attempt to pull any images, as it is expected that those images will be created using the build operation.
`)
}

func (operation *Operation_Pull) Run() {
	operation.log.Info("running pull operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:BUILD")
	operation.nodes.Pull(operation.targets, operation.Registry, operation.force)
}

func (nodes *Nodes) Pull(targets []string, registry string, force bool) {
	for _, target := range nodes.GetTargets(targets) {
		target.node.log = nodes.log.ChildLog("NODE:"+target.node.Name)
		target.node.Pull(registry, force)
	}
}

func (node *Node) Pull(registry string, force bool) bool {
	if node.Do("pull") {

		image, tag := node.GetImageName()

		if !force && node.hasImage() {
			node.log.Info(node.Name+": Node already has an image ["+image+":"+tag+"], so not pulling it again.  You can force this operation if you want to pull this image.")
			return false
		}

		options := docker.PullImageOptions {
			Repository: image,
			OutputStream: os.Stdout,
			RawJSONStream:false,
		}

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

		node.log.Message(node.Name+": Pulling node image ["+image+":"+tag+"] from server ["+options.Registry+"] using auth ["+auth.Username+"] : "+image+":"+tag)
		node.log.DebugObject( LOG_SEVERITY_DEBUG_LOTS, "AUTH USED: ", map[string]string{"Username":auth.Username, "Password":auth.Password, "Email":auth.Email, "ServerAdddress":auth.ServerAddress})

		// ask the docker client to build the image
		err := node.client.PullImage(options, auth)

		if (err!=nil) {
			node.log.Error(node.Name+": Node image not pulled : "+image+" => "+err.Error())
			return false
		} else {
			node.log.Message(node.Name+": Node image pulled: "+image+":"+tag)
			return true
		}

	} else {
		node.log.Info(node.Name+": This node doesn't have an image to pull")
	}
	return true
}
