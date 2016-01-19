package operation

import (
	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type CleanOperation struct {
	log     log.Log
	targets *libs.Targets

	force bool
	wipe bool
	timeout uint
}

func (operation *CleanOperation) Id() string {
	return "clean"
}
func (operation *CleanOperation) Flags(flags []string) bool {
	operation.force = false
	operation.wipe = false
	operation.timeout = 10

	for _, flag := range flags {
		switch flag {
			case "-f":
				fallthrough
			case "--force":
				operation.force = true
			case "-w":
				fallthrough
			case "--wipe":
				operation.wipe = true
			case "-q":
				fallthrough
			case "--quick":
				operation.timeout = 0
		}
	}
	return true
}
func (operation *CleanOperation) Help(topics []string) {
	operation.log.Message(`Operation: Clean

Coach will attempt to Clean any node containers that are active.  Cleaning
means to stop and remove any active containers.  

The operation will also wipe nodes, if the correct flag is passed. Wiping 
means to also remove any built images.

Syntax:
	$/> coach {targets} clean

	$/> coach {targets} clean --wipe   
	  also wipe any built images

	- to eliminate any timeout delays on "docker stop" calls, use this Syntax
	$/> coach {targets} clean --quick

	{targets} what target node instances the operation should process ($/> coach help targets)

`)
}
func (operation *CleanOperation) Run(logger log.Log) bool {
	logger.Message("RUNNING CLEAN OPERATION")
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
			nodeLogger.Info("No node [" + node.MachineName() + "]")
		} else if !node.Can("Clean") {
			nodeLogger.Info("Node doesn't Clean [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Cleaning node")

			if hasInstances {
				for _, id := range instances.InstancesOrder() {
					instance, _ := instances.Instance(id)
					instanceClient := instance.Client()

					if instanceClient.HasContainer() {
						if instanceClient.IsRunning() {
							instanceClient.Stop(logger, operation.force, operation.timeout)
						}
						instanceClient.Remove(logger, operation.force)
						nodeLogger.Message("Cleaning node instance [" + id + "]")						
					} else {
						nodeLogger.Info("Node instance has no container to clean [" + id + "]")					
					}

				}
			}

			if operation.wipe && node.Can("build") {
				nodeClient := node.Client()
				nodeClient.Destroy(nodeLogger, operation.force)
				nodeLogger.Message("Node build cleaned")
			}

		}
	}

	return true
}
