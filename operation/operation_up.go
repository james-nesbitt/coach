package operation

import (
	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type UpOperation struct {
	log     log.Log
	targets *libs.Targets

	force bool
}

func (operation *UpOperation) Id() string {
	return "up"
}
func (operation *UpOperation) Flags(flags []string) bool {
	for _, flag := range flags {
		switch flag {
		case "-f":
			fallthrough
		case "--force":
			operation.force = true
		}
	}
	return true
}
func (operation *UpOperation) Help(topics []string) {
	operation.log.Message(`Operation: Up

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

	{targets} what target node instances the operation should process ($/> coach help targets)

TODO:
	- building images may take a long time, so maybe it should be optional;
	- pulling images may take a long time, so maybe it should be optional;
	- perhaps check if containers exist before creating them;
	- perhaps check if contaienrs are running before starting them.
`)
}
func (operation *UpOperation) Run(logger log.Log) bool {
	logger.Message("RUNNING UP OPERATION")
	logger.Debug(log.VERBOSITY_DEBUG, "Run:Targets", operation.targets.TargetOrder())

	for _, targetID := range operation.targets.TargetOrder() {
		target, targetExists := operation.targets.Target(targetID)
		if !targetExists {
			// this is strange
			logger.Warning("Internal target error, was told to use a target that doesn't exist")
			continue
		}

		node, hasNode := target.Node()
		instances, hasInstances := target.Instances()
		nodeLogger := logger.MakeChild(targetID)

		create := node.Can("create")
		start := node.Can("start")

		if !hasNode {
			nodeLogger.Info("No node [" + node.MachineName() + "]")
		} else if !node.Can("Up") {
			nodeLogger.Info("Node doesn't Up [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Uping node")

			if node.Can("build") {
				nodeLogger.Message("Building node image")
				node.Client().Build(nodeLogger, operation.force)
			}
			if node.Can("pull") {
				nodeLogger.Message("Pulling node image")
				node.Client().Pull(nodeLogger, operation.force)
			}
	
			if hasInstances && (create || start) {
				for _, id := range instances.InstancesOrder() {
					instance, _ := instances.Instance(id)
					instanceClient := instance.Client()

					if create {
						nodeLogger.Message("Creating node instance container : "+id)
						instanceClient.Create(nodeLogger, []string{}, operation.force)
					}
					if start {
						nodeLogger.Message("Starting node instance container : "+id)
						instanceClient.Start(nodeLogger, operation.force)					
					}
				}
			}
		}
	}

	return true
}
