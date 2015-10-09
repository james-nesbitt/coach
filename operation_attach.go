package main

import (
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Attach struct {
	log Log

	nodes Nodes
	targets []string
}
func (operation *Operation_Attach) Flags(flags []string) {

}

func (operation *Operation_Attach) Help(topics []string) {
	operation.log.Note(`Operation: ATTACH

Coach will attempt to attach to an existing container.

SYNTAX:
    $/> coach {target} attach

  {target} what target node instance the operation should process ($/> coach help targets)

ACCESS:
  - This operation processed only nodes with the "start" access.  This excludes build, volume and command containers.

NOTES:
  - The attach operation is really meant to be used to attach to a single instance target, but it multiple instances are targeted, then the operation will attach to each container in sequence.  To target a specific container use an instance style target (for more help checkout $/> coach help target).
`)
}

func (operation *Operation_Attach) Run() {
	operation.nodes.Attach(operation.targets)
}

func (nodes *Nodes) Attach(targets []string) {
	for _, target := range nodes.GetTargets(targets) {
		if target.node.Do("start") {
			for _, instance := range target.instances {
				if instance.HasContainer(true) {
					instance.Attach()
				}
			}
		}
	}
}

func (node *Node) Attach(filters []string) {
	if node.Do("start") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances()
		} else {
			instances = node.FilterInstances(filters)
		}
		for _, instance := range instances {
			if instance.HasContainer(true) {
				instance.Attach()
			}
		}

	}
}

func (instance *Instance) Attach() bool {

	id := instance.GetContainerName()

	// build options for the docker attach operation
	options := docker.AttachToContainerOptions {
		Container:    id,
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream: os.Stdout,

		Logs: true, // Get container logs, sending it to OutputStream.
		Stream: true, // Stream the response?

		Stdin: true, // Attach to stdin, and use InputStream.
		Stdout: true, // Attach to stdout, and use OutputStream.
		Stderr: true,

		//Success chan struct{}

		RawTerminal: instance.Config.Tty, // Use raw terminal? Usually true when the container contains a TTY.
	}

	instance.Node.log.Message("ATTACHING TO INSTANCE CONTAINER ["+id+"]")
	err := instance.Node.client.AttachToContainer( options )
	if err!=nil {
		instance.Node.log.Error("FAILED TO ATTACH TO INSTANCE CONTAINER ["+id+"] =>"+err.Error())
		return false
	} else {
		instance.Node.log.Message("DISCONNECTED FROM INSTANCE CONTAINER ["+id+"]")
		return true
	}

}
