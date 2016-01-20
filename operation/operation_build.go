package operation

import (
	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type BuildOperation struct {
	log     log.Log
	targets *libs.Targets

	force bool
}

func (operation *BuildOperation) Id() string {
	return "build"
}
func (operation *BuildOperation) Flags(flags []string) bool {
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
func (operation *BuildOperation) Help(topics []string) {
	operation.log.Message(`Operation: BUILD

Coach will attempt to build a new docker image, for each target node that has a build setting.

The operation will look for a Build: setting inside the node, and try to find a matching Dockerfile at the suggested path.  The path can be relative to the project root, or absolute.

SYNTAX:
	$/> coach {targets} build

	{targets} what target nodes the operation should process ($/> coach help targets)

ACCESS:
	- this operation processes only nodes with the "build" access.  This includes only nodes with a Build: setting.

NOTES:
- a node that can be built should have a Build: setting, which points to a project path that contains the Dockerfile.
- while {targets} globally can specify particular node instances, that information is ignored for this operation, as images are built for all instances of a node.
`)
}
func (operation *BuildOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: build")
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
		} else if !node.Can("build") {
			nodeLogger.Info("Node doesn't build [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Building node")
			node.Client().Build(nodeLogger, operation.force)
		}
	}

	return true
}
