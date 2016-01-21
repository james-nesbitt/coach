package operation

import (
	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type CreateOperation struct {
	log     log.Log
	targets *libs.Targets

	force bool
}

func (operation *CreateOperation) Id() string {
	return "create"
}
func (operation *CreateOperation) Flags(flags []string) bool {
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

func (operation *CreateOperation) Help(topics []string) {
	operation.log.Message(`Operation: CREATE

Coach will attempt to create any node containers that should be active.

Syntax:
	$/> coach {targets} create

	{targets} what target node instances the operation should process ($/> coach help targets)

Access:
	- only nodes with the "create" access are processed.  This excludes build and command nodes
`)
}
func (operation *CreateOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: create")
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
		} else if !node.Can("create") {
			nodeLogger.Info("Node doesn't create [" + node.MachineName() + ":" + node.Type() + "]")
		} else if !hasInstances {
			nodeLogger.Info("No valid instances specified in target list [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Creating instance containers")
			for _, id := range instances.InstancesOrder() {
				instance, _ := instances.Instance(id)
				instance.Client().Create(logger, []string{}, operation.force)
			}
		}
	}

	return true
}
