package operation

import (
	"strconv"

	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type ScaleOperation struct {
	log     log.Log
	targets *libs.Targets

	force         bool
	scale         int
	timeout       uint
	removeStopped bool
}

func (operation *ScaleOperation) Id() string {
	return "scale"
}
func (operation *ScaleOperation) Flags(flags []string) bool {
	// default scale value
	operation.scale = 1
	operation.force = true
	operation.timeout = 5
	operation.removeStopped = true

	if len(flags) > 0 {
		if flags[0] == "up" {
			operation.scale = 1
		} else if flags[0] == "down" {
			operation.scale = -1
		}
	}
	return true
}
func (operation *ScaleOperation) Help(topics []string) {
	operation.log.Message(`Operation: Scale

Coach will attempt to scale up or down, the number of running instances on a scaled node

Syntax:
	$/> coach {targets} scale

	{targets} what target node instances the operation should process ($/> coach help targets)

Access:
	- only scaled type nodes with the "start" access are processed.  This effectively limits it to service type nodes.
`)
}
func (operation *ScaleOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: scale")
	logger.Debug(log.VERBOSITY_DEBUG, "Run:Targets", operation.targets.TargetOrder())

	if operation.scale == 0 {
		operation.log.Warning("scale operation was told to scale to 0")
		return false
	}

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
		} else if !node.Can("scale") {
			nodeLogger.Info("Node doesn't Scale [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Scaling node " + node.Id())

			if operation.scale > 0 {
				count := operation.ScaleUpNumber(nodeLogger, node.Instances(), operation.scale)

				if count == 0 {
					nodeLogger.Warning("Scale operation could not scale up any new instances of node")
				} else if count < operation.scale {
					nodeLogger.Warning("Scale operation could not scale up all of the requested instances of node. " + strconv.FormatInt(int64(count+1), 10) + " started.")
				} else {
					nodeLogger.Message("Scale operation scaled up " + strconv.FormatInt(int64(count), 10) + " instances")
				}

			} else {
				count := operation.ScaleDownNumber(nodeLogger, node.Instances(), -operation.scale)

				if count == 0 {
					nodeLogger.Warning("Scale operation could not scale down any new instances of node")
				} else if count < (-operation.scale) {
					nodeLogger.Warning("Scale operation could not scale down all of the requested instances of node. " + strconv.FormatInt(int64(count+1), 10) + " stopped.")
				} else {
					nodeLogger.Message("Scale operation scaled down " + strconv.FormatInt(int64(count), 10) + " instances")
				}
			}

		}
	}
	return true
}

// scale a node up a certain number of instances
func (operation *ScaleOperation) ScaleUpNumber(logger log.Log, instances libs.Instances, number int) int {
	count := 0
	instancesOrder := instances.InstancesOrder()

InstanceScaleReturn:
	for _, instanceId := range instancesOrder {
		if instance, ok := instances.Instance(instanceId); ok {
			client := instance.Client()

			if client.IsRunning() {
				continue InstanceScaleReturn
			} else if !client.HasContainer() {
				// create a new container for this instance
				client.Create(logger, []string{}, false)
			}

			logger.Info("Node Scaling up. Starting instance :" + instanceId)
			client.Start(logger, false)

			count++
			if count >= number {
				return count
			}
		}
	}

	return count
}

// scale a node down a certain number of instances
func (operation *ScaleOperation) ScaleDownNumber(logger log.Log, instances libs.Instances, number int) int {

	count := 0
	instancesOrder := []string{}
	for _, instanceId := range instances.InstancesOrder() {
		instancesOrder = append([]string{instanceId}, instancesOrder...)
	}

InstanceScaleReturn:
	for _, instanceId := range instancesOrder {
		if instance, ok := instances.Instance(instanceId); ok {
			client := instance.Client()

			if !client.IsRunning() {
				continue InstanceScaleReturn
			}

			logger.Info("Node Scaling down. Stopping instance :" + instanceId)
			client.Stop(logger, operation.force, operation.timeout)

			if operation.removeStopped {
				client.Remove(logger, operation.force)
			}

			count++
			if count >= number {
				return count
			}
		}
	}

	return count
}
