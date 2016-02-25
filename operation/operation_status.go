package operation

import (
	"strings"

	"github.com/james-nesbitt/coach/libs"
	"github.com/james-nesbitt/coach/log"
)

type StatusOperation struct {
	log     log.Log
	targets *libs.Targets
}

func (operation *StatusOperation) Id() string {
	return "Status"
}
func (operation *StatusOperation) Flags(flags []string) bool {
	return true
}
func (operation *StatusOperation) Help(flags []string) {
	operation.log.Message(`Operation: Status

Coach will attempt to provide project Status by investigating target images and containers.

SYNTAX:
	$/> coach {targets} Status

	{targets} what target nodes the operation should process ($/> coach help targets)

`)
}
func (operation *StatusOperation) Run(logger log.Log) bool {
	logger.Message("RUNNING Status OPERATION")

	logger.Debug(log.VERBOSITY_DEBUG, "Run:Targets", operation.targets.TargetOrder())
	for _, targetID := range operation.targets.TargetOrder() {
		target, targetExists := operation.targets.Target(targetID)
		node, hasNode := target.Node()
		instances, hasInstances := target.Instances()

		if !targetExists {
			// this is strange
			logger.Warning("Internal target error, was told to use a target that doesn't exist")
			continue
		}

		nodeLogger := logger.MakeChild(targetID)
		status := []string{}


		if hasNode {
			status = append(status, node.Status(nodeLogger)...)
		} else {
			status = append(status, "No node for target")
		}
		if hasInstances {
			status = append(status, instances.Status(nodeLogger)...)
		} else {
			status = append(status, "No instances for target")
		}

		nodeLogger.Message( "["+strings.Join(status, "][")+"]" )
	}

	return true
}
