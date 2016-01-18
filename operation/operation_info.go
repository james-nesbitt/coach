package operation

import (
	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type InfoOperation struct {
	log     log.Log
	targets *libs.Targets
}

func (operation *InfoOperation) Id() string {
	return "info"
}
func (operation *InfoOperation) Flags(flags []string) bool {

	return true
}
func (operation *InfoOperation) Run(logger log.Log) bool {
	logger.Message("RUNNING INFO OPERATION")

	logger.Debug(log.VERBOSITY_DEBUG, "Run:Targets", operation.targets.TargetOrder())
	for _, targetID := range operation.targets.TargetOrder() {
		target, targetExists := operation.targets.Target(targetID)
		node, hasNode := target.Node()
		_, hasInstances := target.Instances()

		if !targetExists {
			// this is strange
			logger.Warning("Internal target error, was told to use a target that doesn't exist")
			continue
		}

		nodeLogger := logger.MakeChild(targetID)

		if hasNode {
			nodeLogger.Message(targetID + " Information")
			node.Client().NodeInfo(nodeLogger)
		} else {
			nodeLogger.Message("No node [" + node.MachineName() + "]")
		}
		if hasInstances {
			node.Instances().Client().InstancesInfo(nodeLogger)
		} else {
			nodeLogger.Message("|-- No instances [" + node.MachineName() + "]")
		}
	}

	return true
}
