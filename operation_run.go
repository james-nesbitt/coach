package main

import (
	"strings"
)

type Operation_Run struct {
	log Log

	nodes Nodes
	targets []string

	cmd []string
	instance string
}
func (operation *Operation_Run) Flags(flags []string) {
	if len(flags)>0 && strings.HasPrefix(flags[0], "->") {
		operation.instance = string(flags[0][1:])

		if len(flags)>1 {
			flags = flags[1:]
		} else {
			flags = []string{}
		}
	}

	operation.cmd = flags
}

func (operation *Operation_Run) Help(topics []string) {
	operation.log.Note(`Operation: RUN

Coach will attempt a single command run on a node container.

The run operation follows the following steps:
- creates a new container using a new command (read from command line)
- starts that container, output stdout and stderr
- removes the started container

The process is ideal for running single commands in volatile containers, which can disappear after execution.

Containers can be persistant, but such containers are not as usefull, as the command cannot be changed.  In most cases, volatility can still work, as long as persistant file and folder maps are used to keep volatile information.
`)
}

func (operation *Operation_Run) Run() {
	operation.log.Message("running run operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)

	for _, target := range operation.nodes.GetTargets(operation.targets, true) {
		target.node.log = operation.nodes.log.ChildLog("NODE:"+target.node.Name)
		target.node.Run(operation.instance, operation.cmd)
	}
}

func (node *Node) Run(instanceid string, cmd []string) bool {
	if node.Do("run") {

		var instance *Instance
		var persistant bool

		switch node.InstanceType {
			case "temporary":
				if instanceid=="" {
					instanceid = "run" // perhaps we should randomly generate this to allow for persistance run containers
				}
				// create a new temporary instance for a node
				node.AddTemporaryInstance(instanceid)
				instance = node.GetInstance(instanceid)
				persistant = false
			case "single":
				instanceid = "single"
				fallthrough
			default:
				instance = node.GetInstance(instanceid)
				persistant = true
		}

		if instance!=nil {
			return instance.Run(cmd, persistant)
		}

	}
	return false
}

func (instance *Instance) Run(cmd []string, persistant bool) bool {

	// Set up some additional settings for TTY commands
	if instance.Config.Tty==true {

		// set a default hostname to make a prettier prompt
		if instance.Config.Hostname=="" {
			instance.Config.Hostname = instance.GetContainerName()
		}

		// make sure that all tty runs have openstdin
		instance.Config.OpenStdin=true

	}

	instance.Config.AttachStdin = true
	instance.Config.AttachStdout = true
	instance.Config.AttachStderr = true

	// 1. get the container for the instance (create it if needed)
	hasContainer := instance.HasContainer(false)
	if !hasContainer {
		hasContainer = instance.Create(cmd, false)
	}
	if hasContainer {

	// 3. start the container (set up a remove)
		ok := instance.Start(false)

	// 4. attach to the container
		if ok {
			if !persistant {
				// 5. remove the container (if not instructed to keep it)
				defer instance.Remove(true)
			}
			instance.Attach()
			return true
		} else {
			instance.Node.log.Error("Could not start RUN container")
			return false
		}

	} else {
		instance.Node.log.Error("Could not create RUN container")
	}

	return false
}
