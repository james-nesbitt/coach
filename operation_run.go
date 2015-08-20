package main

import (
	"strings"
)

type Operation_Run struct {
	log Log

	Nodes Nodes
	Targets []string

	cmd []string
	instance string
}
func (operation *Operation_Run) Flags(flags []string) {
	if len(flags)>0 && strings.HasPrefix(flags[0], "@") {
		operation.instance = string(flags[0][1:])

		if len(flags)>1 {
			flags = flags[1:]
		} else {
			flags = []string{}
		}
	}

	operation.cmd = flags
}

func (operation *Operation_Run) Help() {

}

func (operation *Operation_Run) Run() {
	operation.log.Message("running run operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.Targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:RUN")
	operation.Nodes.Run(operation.Targets, operation.instance, operation.cmd)
}

func (nodes *Nodes) Run(targets []string, instance string, cmd []string) {
	for _, target := range nodes.GetTargets(targets) {
		target.log = nodes.log.ChildLog("NODE:"+target.Name)
		target.Run(instance, cmd)
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
