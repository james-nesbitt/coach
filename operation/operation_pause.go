package operation

import (
	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type PauseOperation struct {
	log     log.Log
	targets *libs.Targets
}

func (operation *PauseOperation) Id() string {
	return "pause"
}
func (operation *PauseOperation) Flags(flags []string) bool {
	return true
}
func (operation *PauseOperation) Help(topics []string) {
	operation.log.Message(`Operation: Pause


Coach will attempt to pause any target containers.

SYNTAX:
	$/> coach {targets} pause

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
	- This operation processed only nodes with the "start" access.  This excludes build, volume and command containers

`)
}
func (operation *PauseOperation) Run(logger log.Log) bool {
	logger.Message("RUNNING PAUSE OPERATION")
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
		} else if !node.Can("pause") {
			nodeLogger.Info("Node doesn't Pause [" + node.MachineName() + ":" + node.Type() + "]")
		} else if !hasInstances {
			nodeLogger.Info("No valid instances specified in target list [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Pausing instance containers")
			for _, id := range instances.InstancesOrder() {
				instance, _ := instances.Instance(id)
				instance.Client().Pause(logger)
			}
		}
	}

	return true
}
