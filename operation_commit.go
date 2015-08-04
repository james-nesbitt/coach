package main

import (
	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Commit struct {
	log Log

	Nodes Nodes
	Targets []string
}
func (operation *Operation_Commit) Flags(flags []string) {

}
func (operation *Operation_Commit) Run() {
	operation.Nodes.Commit(operation.Targets, "/", map[string]string{"single":"latest"}, "")
}

func (nodes *Nodes) Commit(targets []string, repo string, instanceTags map[string]string, message string) {
	for _, target := range nodes.GetTargets(targets) {
		for instance, tag := range instanceTags {
			target.Commit(repo, instance, tag, message)
		}
	}
}

func (node *Node) Commit(repo string, instance string, tag string, message string) {
	if node.Do("commit") {

		for _, instance := range node.FilterInstances([]string{instance}, false) {
			if instance.HasContainer(false) {
				instance.Commit(repo, tag, message)
			}
		}

	}
}

func (instance *Instance) Commit(repo string, tag string, message string) bool {

	id := instance.GetContainerName()
	config := instance.Config

	if (config.Image=="") {
		config.Image = instance.Node.GetImageName()
	}

	options := docker.CommitContainerOptions{
		Container: id,
		Repository: repo,
		Tag: tag,
		Run: &config,
	}

	if message!="" {
		options.Message = message
	}
	if instance.Node.conf.Author!="" {
		options.Author = instance.Node.conf.Author
	}

	_, err := instance.Node.client.CommitContainer( options )
	if err!=nil {
		instance.Node.log.Warning("Failed to commit container changes to an image ["+instance.Node.Name+":"+id+"] : "+tag)
		return false
	} else {
		instance.Node.log.Message("Committed container changes to an image ["+instance.Node.Name+":"+id+"] : "+tag)
		return true
	}
}
