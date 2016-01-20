package operation

import (
	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type UnpauseOperation struct {
	log     log.Log
	targets *libs.Targets
}

func (operation *UnpauseOperation) Id() string {
	return "unpause"
}
func (operation *UnpauseOperation) Flags(flags []string) bool {
	return true
}
func (operation *UnpauseOperation) Help(topics []string) {
	operation.log.Message(`Operation: UnPause

Coach will attempt to unpause target node containers.

SYNTAX:
	$/> coach {targets} unpause

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
	- This operation processed only nodes with the "start" access.  This excludes build, volume and command containers.
`)
}
func (operation *UnpauseOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: unpause")
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
		} else if !node.Can("npause") {
			nodeLogger.Info("Node doesn't unpause [" + node.MachineName() + ":" + node.Type() + "]")
		} else if !hasInstances {
			nodeLogger.Info("No valid instances specified in target list [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("UnPausing instance containers")

			if !instances.IsFiltered() {
				nodeLogger.Message("Switching to using all instances")
				instances.UseAll()
			}

			for _, id := range instances.InstancesOrder() {
				instance, _ := instances.Instance(id)

				if instance.IsRunning() {
					instance.Client().Unpause(logger)
				}
			}
		}
	}

	return true
}
