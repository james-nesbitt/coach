package operation

import (
	"github.com/james-nesbitt/coach/libs"
	"github.com/james-nesbitt/coach/log"
)

type StopOperation struct {
	log     log.Log
	targets *libs.Targets

	force   bool
	timeout uint
}

func (operation *StopOperation) Id() string {
	return "stop"
}

func (operation *StopOperation) Flags(flags []string) bool {
	operation.timeout = 10

	for index := 0; index < len(flags); index++ {
		flag := flags[index]

		switch flag {
		case "-q":
			fallthrough
		case "--quick":
			operation.timeout = 1

		}
	}
	return true
}

func (operation *StopOperation) Help(topics []string) {
	operation.log.Message(`Operation: Stop

Coach will attempt to stop target node containers.

SYNTAX:
	$/> coach {targets} stop

	$/> coach {targets} stop --quick
	- makes docker stop the containers with --time=1

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
	- This operation processed only nodes with the "start" access.  This excludes build, volume and command containers.

`)
}
func (operation *StopOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: stop")
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
		} else if !node.Can("Stop") {
			nodeLogger.Info("Node doesn't Stop [" + node.MachineName() + ":" + node.Type() + "]")
		} else if !hasInstances {
			nodeLogger.Info("No valid instances specified in target list [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Stopping instance containers")

			if !instances.IsFiltered() {
				nodeLogger.Info("Switching to using all instances")
				instances.UseAll()
			}

			for _, id := range instances.InstancesOrder() {
				instance, _ := instances.Instance(id)

				if instance.IsRunning() {
					instance.Client().Stop(logger, operation.force, operation.timeout)
				} else {
					nodeLogger.Info("Instance [" + id + "] is not running")
				}
			}
		}
	}

	return true
}
