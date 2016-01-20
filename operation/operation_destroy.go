package operation

import (
	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type DestroyOperation struct {
	log     log.Log
	targets *libs.Targets

	force bool
}

func (operation *DestroyOperation) Id() string {
	return "Destroy"
}
func (operation *DestroyOperation) Flags(flags []string) bool {
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
func (operation *DestroyOperation) Help(topics []string) {
	operation.log.Message(`Operation: Destroy

Coach will attempt to remove any built images for target nodes.  A node is
considered "buildable" if it has a Build: definition.

SYNTAX:
	$/> coach {targets} destroy

	{targets} what target nodes the operation should process ($/> coach help targets)

ACCESS:
	- this opertion will only process nodes with "build" access.  This includes only nodes with the Build: settings declared.

NOTE:
	- Coach will try not to remove an image for a node that does not build, which is the common case for nodes with an image, but no build setting.  This prevents deleting shared images that are not build targets.
`)
}
func (operation *DestroyOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: destroy")
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
		} else if !node.Can("Destroy") {
			nodeLogger.Info("Node doesn't Destroy [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Destroying node")
			node.Client().Destroy(nodeLogger, operation.force)
		}
	}

	return true
}
