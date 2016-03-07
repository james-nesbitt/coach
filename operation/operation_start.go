package operation

import (
	"github.com/james-nesbitt/coach/libs"
	"github.com/james-nesbitt/coach/log"
)

type StartOperation struct {
	log     log.Log
	targets *libs.Targets

	force bool
}

func (operation *StartOperation) Id() string {
	return "start"
}
func (operation *StartOperation) Flags(flags []string) bool {
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

func (operation *StartOperation) Help(topics []string) {
	operation.log.Message(`Operation: Start

Coach will attempt to start target node containers.

SYNTAX:
	$/> coach {targets} start

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
	- This operation processed only nodes with the "start" access.  This excludes build, volume and command containers.
`)
}
func (operation *StartOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: start")
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

		if !hasNode {
			nodeLogger.Warning("No node [" + node.MachineName() + "]")
		} else if !node.Can("start") {
			nodeLogger.Info("Node doesn't Start [" + node.MachineName() + ":" + node.Type() + "]")
		} else if !hasInstances {
			nodeLogger.Info("No valid instances specified in target list [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Starting instance containers")
			for _, id := range instances.InstancesOrder() {
				instance, _ := instances.Instance(id)
				instance.Client().Start(logger, operation.force)
			}
		}
	}

	return true
}
