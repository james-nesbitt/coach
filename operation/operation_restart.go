package operation

import (
	"github.com/james-nesbitt/coach/libs"
	"github.com/james-nesbitt/coach/log"
)

type RestartOperation struct {
	log     log.Log
	targets *libs.Targets

	force bool
	timeout uint
}

func (operation *RestartOperation) Id() string {
	return "Restart"
}
func (operation *RestartOperation) Flags(flags []string) bool {
	for _, flag := range flags {
		switch flag {
		case "-f":
			fallthrough
		case "--force":
			operation.force = true


		case "-q":
			fallthrough
		case "--quick":
			operation.timeout = 1	
		}			
	}
	return true
}

func (operation *RestartOperation) Help(topics []string) {
	operation.log.Message(`Operation: Restart

Coach will attempt to Restart target node containers.

SYNTAX:
	$/> coach {targets} Restart

	$/> coach {targets} Restart--quick
	- makes docker stop the containers with --time=1	

	{targets} what target node instances the operation should process ($/> coach help targets)

ACCESS:
	- This operation processed only nodes with the "Restart" access.  This excludes build, volume and command containers.
`)
}
func (operation *RestartOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: Restart")
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
		} else if !node.Can("start") {
			nodeLogger.Info("Node doesn't Restart [" + node.MachineName() + ":" + node.Type() + "]")
		} else if !hasInstances {
			nodeLogger.Info("No valid instances specified in target list [" + node.MachineName() + "]")
		} else {
			nodeLogger.Info("Restarting instance containers")

			for _, id := range instances.InstancesOrder() {
				instance, _ := instances.Instance(id)
				instanceClient := instance.Client()

				// Stop
				if instanceClient.IsRunning() {
					nodeLogger.Info("Stopping instance: "+id)
					instanceClient.Stop(nodeLogger, operation.force, operation.timeout)
				}

				// Start
				if !instanceClient.IsRunning() {
					nodeLogger.Info("Starting instance: "+id)
					instanceClient.Start(nodeLogger, operation.force)
				}

			}
		}
	}

	return true
}
