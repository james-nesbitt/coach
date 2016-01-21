package operation

import (
	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type RunOperation struct {
	log     log.Log
	targets *libs.Targets

	cmd      []string
	instance string

	persistant bool
}

func (operation *RunOperation) Id() string {
	return "Run"
}
func (operation *RunOperation) Flags(flags []string) bool {
	operation.cmd = flags
	return true
}
func (operation *RunOperation) Help(topics []string) {
	operation.log.Message(`Operation: Run

Coach will attempt a single command run on a node container.

The run operation follows the following steps:
- creates a new container using a new command (read from command line)
- starts that container, output stdout and stderr
- removes the started container

The process is ideal for running single commands in volatile containers, which can disappear after execution.

SYNTAX:
	$/> coach {target} run {cmd}

	{target} what target node instance the operation should process ($/> coach help targets)
	{cmd} a list of flags to pass into the container.  These can be flags added passed to the container entrypoint, or full command replacement.

NOTE:
-	 Containers can be persistant, but such containers are generally not usefull, as the container command cannot be changed.  In most cases, command container volatility can still work, as long as persistant file and folder binds/maps are used to keep volatile information outside of the container.

TODO:
	- Allow overriding of a container entrypoint via a flag?
`)
}
func (operation *RunOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: run")
	logger.Debug(log.VERBOSITY_DEBUG, "Run:Targets", operation.targets.TargetOrder())

	for _, targetID := range operation.targets.TargetOrder() {
		target, targetExists := operation.targets.Target(targetID)
		if !targetExists {
			// this is strange
			logger.Warning("Internal target error, was told to use a target that doesn't exist")
			continue
		}

		node, hasNode := target.Node()
		instances, _ := target.Instances()
		nodeLogger := logger.MakeChild(targetID)

		if !hasNode {
			nodeLogger.Info("No node [" + node.MachineName() + "]")
		} else if !node.Can("run") {
			nodeLogger.Info("Node doesn't Run [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Runing node")

			instanceIds := instances.InstancesOrder()
			if len(instanceIds) == 0 {
				instanceIds = []string{""}
			}

			for _, id := range instanceIds {
				if instance, ok := instances.Instance(id); ok {
					instance.Client().Run(logger, false, operation.cmd)
				}
			}
		}
	}

	return true
}
