package main

import (
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

type Operation_Attach struct {
	log Log

	Nodes Nodes
	Targets []string

}
func (operation *Operation_Attach) Flags(flags []string) {

}

func (operation *Operation_Attach) Help() {

}

func (operation *Operation_Attach) Run() {
	operation.Nodes.Attach(operation.Targets)
}

func (nodes *Nodes) Attach(targets []string) {
	for _, target := range nodes.GetTargets(targets) {
		target.Attach([]string{})
	}
}

func (node *Node) Attach(filters []string) {
	if node.Do("start") {

		var instances []*Instance

		if len(filters)==0 {
			instances = node.GetInstances(false)
		} else {
			instances = node.FilterInstances(filters, false)
		}
		for _, instance := range instances {
			if instance.active {
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

		RawTerminal: false, // Use raw terminal? Usually true when the container contains a TTY.
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
