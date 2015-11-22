package main

import (
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Commit struct {
	log Log

	nodes Nodes
	targets []string

	repo string

	tag string
	message string
}
func (operation *Operation_Commit) Flags(flags []string) {

	for index:=0; index<len(flags); index++ {
		flag:= flags[index]

		switch flag {
			case "-t":
				fallthrough
			case "--tag":
				if !strings.HasPrefix(flags[index+1], "-") {
					operation.tag = flags[index+1]
					index++
				}
			case "-r":
				fallthrough
			case "--repo":
				if !strings.HasPrefix(flags[index+1], "-") {
					operation.repo = flags[index+1]
					index++
				}
			case "-m":
				fallthrough
			case "--message":
				if !strings.HasPrefix(flags[index+1], "-") {
					operation.repo = flags[index+1]
					index++
				}
		}
	}

}

func (operation *Operation_Commit) Help(topics []string) {
	operation.log.Note(`Operation: COMMIT

Coach will attempt to commit a container to it's image.

SYNTAX:
    $/> coach {targets} commit [--tag {tag}] [--repo {repo}] [--message "{message}"]

  {targets} what target node instances the operation should process ($/> coach help targets)
  --tag "{tag}" : what image tag to use (default: "latest")
  --repo "{repo}" : what image repository to commit to (default: local)
  --message "{message}" : what commit message to use

ACCESS:
  - only nodes with the "commit" access are processed.  This excludes build nodes

`)
}

func (operation *Operation_Commit) Run() {
	if operation.tag=="" {
		operation.tag = "latest"
	}
	if operation.repo=="" {
		operation.repo = ""
	}

	operation.nodes.Commit(operation.targets, operation.repo, operation.tag, "")
}

func (nodes *Nodes) Commit(targets []string, repo string, tag string, message string) {
	for _, target := range nodes.GetTargets(targets) {
		if target.node.Do("commit") {
			for _, instance := range target.instances {
				if instance.HasContainer(false) {
					instance.Commit(repo, tag, message)
				}
			}
		}
	}
}

func (node *Node) Commit(repo string, instance string, tag string, message string) {
	if node.Do("commit") {

		for _, instance := range node.FilterInstances([]string{instance}) {
			if instance.HasContainer(false) {
				instance.Commit(repo, tag, message)
			}
		}

	}
}

func (instance *Instance) Commit(repo string, tag string, message string) bool {

	id := instance.GetContainerName()
	config := instance.Config

	if repo=="" {
		repo, _ = instance.Node.GetImageName()
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
