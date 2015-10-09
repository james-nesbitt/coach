package main

import (
	"strconv"
	"strings"
)

type Operation_Scale struct {
	log Log

	nodes Nodes
	targets []string
}
func (operation *Operation_Scale) Flags(flags []string) {

}

func (operation *Operation_Scale) Help(topics []string) {
	operation.log.Note(`Operation: CREATE

Coach will attempt to scale up or down, the number of running instances on a scaled node

Syntax:
    $/> coach {targets} scale

  {targets} what target node instances the operation should process ($/> coach help targets)

Access:
  - only scaled type nodes with the "start" access are processed.  This effectively limits it to service type nodes.
`)
}

func (operation *Operation_Scale) Run() {
	operation.log.Message("running scale operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)

	TargetScaleReturn:
	for _, target := range operation.nodes.GetTargets(operation.targets) {

		InstanceScaleReturn:
		for i := 0; i < len(target.node.InstanceMap); i++ {
			if instance, ok := target.node.InstanceMap[strconv.FormatInt(int64(i+1), 10)]; ok {

				container, hasContainer := instance.GetContainer(false)

				if hasContainer {
					if strings.Contains(container.Status, "Up") {
						operation.log.Debug(LOG_SEVERITY_DEBUG_LOTS, "Node instance already has a running container :"+container.ID)
						continue InstanceScaleReturn
					}
				} else {
					// create a new container for this instance
					instance.Create([]string{}, false)
				}

				operation.log.Message("Scaling up instance")
				instance.Start(false)
				continue TargetScaleReturn

			}
// operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "SCALE NODE TEST INSTANCE:", *target.node.InstanceMap[strconv.FormatInt(int64(i+1), 10)])
		}

		operation.log.Warning("Could not find any node instance to scale up to :"+target.node.Name)

	}

	cache.refresh(false, true)
}
