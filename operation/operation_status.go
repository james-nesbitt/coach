package operation

import (
	"strconv"
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
			status = append(status, operation.NodeStatus(nodeLogger, node)...)
		} else {
			status = append(status, "No node for target")
		}
		if hasInstances {
			status = append(status, operation.InstancesStatus(nodeLogger, instances)...)
		} else {
			status = append(status, "No instances for target")
		}

		nodeLogger.Message("[" + strings.Join(status, "][") + "]")
	}

	return true
}

func (operation *StatusOperation) NodeStatus(logger log.Log, node libs.Node) []string {
	status := []string{}

	if node.Client().HasImage() {
		status = append(status, "Image:good")
	} else {
		status = append(status, "Image:No Image")
	}

	return status
}

func (operation *StatusOperation) InstancesStatus(logger log.Log, instances libs.FilterableInstances) []string {
	status := []string{}

	allCount := 0
	runningCount := 0
	for _, id := range instances.InstancesOrder() {
		instance, _ := instances.Instance(id)
		instanceClient := instance.Client()
		if instanceClient.IsRunning() {
			allCount++
			runningCount++
		} else if instanceClient.HasContainer() {
			allCount++
		}
	}

	if allCount > 0 {
		status = append(status, "Containers:"+strconv.Itoa(runningCount)+"/"+strconv.Itoa(allCount))
	} else {
		status = append(status, "Containers:NONE")
	}

	return status
}
