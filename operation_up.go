package main

type Operation_Up struct {
	log Log

	nodes Nodes
	targets[] string

	force bool
}
func(operation * Operation_Up) Flags(flags[] string) {
	operation.force = false

	for _, flag := range flags {
		switch flag {
			case "-f":
				fallthrough
			case "--force":
				operation.force = true
		}
	}

}

func(operation * Operation_Up) Help(topics[] string) {
	operation.log.Note(`Operation: UP

Coach will attempt to get all nodes operational and running.

The run operation follows the following steps:
- build any images for targets with build settings
- pull any required images for targets without build settings
- create new containers for all targets
- start any containers that should be started

This operation is used to allow users to take a coach project
from beginning to fully operational, with a single command.

SYNTAX:
    $/> coach {target} up

TODO:
- building images may take a long time, so maybe it should be optional;
- pulling images may take a long time, so maybe it should be optional;
- perhaps check if containers exist before creating them;
- perhaps check if contaienrs are running before starting them.
`)
}

func (operation * Operation_Up) Run() {
	operation.log.Message("running run operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)

	targets := operation.nodes.GetTargets(operation.targets)

	for _, target := range targets {
		target.node.log = operation.nodes.log.ChildLog("NODE:" + target.node.Name)

		if target.node.Do("build") {
			target.node.Build(operation.force)
		}
		if target.node.Do("pull") {
			target.node.Pull("https://index.docker.io/v1/", operation.force)
		}
	}
	for _, target := range targets {
		for _, instance := range target.instances {
			if instance.Node.Do("create") {
				instance.Create([] string {}, operation.force)
			}
		}
	}
	for _, target := range targets {
		for _, instance := range target.instances {
			if instance.Node.Do("start") {
				instance.Start(operation.force)
			}
		}
	}
}