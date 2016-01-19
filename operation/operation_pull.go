package operation

import (
	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type PullOperation struct {
	log     log.Log
	targets *libs.Targets

	force bool
}

func (operation *PullOperation) Id() string {
	return "pull"
}
func (operation *PullOperation) Flags(flags []string) bool {
	operation.force = false

	for index := 0; index < len(flags); index++ {
		flag := flags[index]

		switch flag {
		case "-f":
			fallthrough
		case "--force":
			operation.force = true
		}
	}
	return true
}
func (operation *PullOperation) Help(topics []string) {
	operation.log.Message(`Operation: PULL

Coach will attempt to pull any node images, for nodes that have no build settings.


SYNTAX:
	$/> coach {targets} pull

	{targets} what target nodes the operation should process ($/> coach help targets)

ACCESS:
	- this operation processes only nodes with the "pull" access.  This includes only nodes without a Build: setting, but a Config: Image: setting.

NOTES:
	- Nodes that have build settings will not attempt to pull any images, as it is expected that those images will be created using the build operation.
`)
}
func (operation *PullOperation) Run(logger log.Log) bool {
	logger.Message("RUNNING PULL OPERATION")

	logger.Debug(log.VERBOSITY_DEBUG, "Run:Targets", operation.targets.TargetOrder())
	for _, targetID := range operation.targets.TargetOrder() {
		target, targetExists := operation.targets.Target(targetID)
		if !targetExists {
			// this is strange
			logger.Warning("Internal target error, was told to use a target that doesn't exist")
			continue
		}

		node, hasNode := target.Node()
		nodeLogger := logger.MakeChild(targetID)

		if !hasNode {
			nodeLogger.Info("No node [" + node.MachineName() + "]")
		} else if !node.Can("pull") {
			nodeLogger.Info("Node doesn't pull [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Pulling node")
			node.Client().Pull(nodeLogger, operation.force)
		}
	}

	return true
}
