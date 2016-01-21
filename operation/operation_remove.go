package operation

import (
	"github.com/james-nesbitt/coach/libs"
	"github.com/james-nesbitt/coach/log"
)

type RemoveOperation struct {
	log     log.Log
	targets *libs.Targets

	force bool
}

func (operation *RemoveOperation) Id() string {
	return "remove"
}
func (operation *RemoveOperation) Flags(flags []string) bool {
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

func (operation *RemoveOperation) Help(topics []string) {
	operation.log.Message(`Operation: REMOVE

Coach will attempt to remove all target node containers.

SYNTAX:
	$/> coach {targets} remove

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
	- only nodes with the "create" access are processed.  This excludes build and command nodes
`)
}
func (operation *RemoveOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: remove")
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
		} else if !node.Can("remove") {
			nodeLogger.Info("Node doesn't remove [" + node.MachineName() + ":" + node.Type() + "]")
		} else if !hasInstances {
			nodeLogger.Info("No valid instances specified in target list [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Removing instance containers")

			if !instances.IsFiltered() {
				nodeLogger.Info("Switching to using all instances")
				instances.UseAll()
			}

			for _, id := range instances.InstancesOrder() {
				instance, _ := instances.Instance(id)
				instance.Client().Remove(logger, operation.force)
			}
		}
	}

	return true
}
