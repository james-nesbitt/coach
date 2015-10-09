package main

import (
	"strconv"
	"strings"
)

type Operation_Scale struct {
	log Log

	nodes Nodes
	targets []string

	scale int
}
func (operation *Operation_Scale) Flags(flags []string) {

	//@TODO GET SCALE FROM FLAG
	operation.scale = 1

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

	if operation.scale==0 {
		operation.log.Warning("scale operation was told to scale to 0")
		return
	}

	TargetScaleReturn:
	for _, target := range operation.nodes.GetTargets(operation.targets) {

		if target.node.InstanceType!="scaled" {
			operation.log.Warning("Tried to scale non-scaleable node :"+target.node.Name)
			continue TargetScaleReturn
		}

		if operation.scale>0 {
			count := target.node.ScaleUpNumber(operation.scale)
			
			if count==0 {
				target.node.log.Warning("Scale operation could not scale up any new instances of node :"+target.node.Name)
			} else if count<operation.scale {
				target.node.log.Warning("Scale operation could not scale up all of the requested instances of node :"+target.node.Name)
			} else {
				target.node.log.Warning("Scale operation scaled up all requested node instances :"+target.node.Name)
			}

		} else {
			count := target.node.ScaleDOwnNumber(-operation.scale)

			if count==0 {
				target.node.log.Warning("Scale operation could not scale down any new instances of node :"+target.node.Name)
			} else if count<(-operation.scale) {
				target.node.log.Warning("Scale operation could not scale down all of the requested instances of node :"+target.node.Name)
			} else {
				target.node.log.Warning("Scale operation scaled down all requested node instances :"+target.node.Name)
			}
		}

	}

	cache.refresh(false, true)
}

func (node *Node) ScaleUpNumber(number int) int {

	count := 0

	InstanceScaleReturn:
	for i := 0; i < len(node.InstanceMap); i++ {
		if instance, ok := node.InstanceMap[strconv.FormatInt(int64(i+1), 10)]; ok {

			container, hasContainer := instance.GetContainer(false)

			if hasContainer {
				if strings.Contains(container.Status, "Up") {
					continue InstanceScaleReturn
				}
			} else {
				// create a new container for this instance
				instance.Create([]string{}, false)
			}

			node.log.Message("Node Scaling up. Starting instance :"+instance.Name)
			instance.Start(false)

			count++
			if count >= number {
				return count
			}
		}
// operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "SCALE NODE TEST INSTANCE:", *target.node.InstanceMap[strconv.FormatInt(int64(i+1), 10)])
	}
	
	return count
}
func (node *Node) ScaleDOwnNumber(number int) int {

	count := 0

	InstanceScaleReturn:
	for i := len(node.InstanceMap); i>=0; i-- {
		if instance, ok := node.InstanceMap[strconv.FormatInt(int64(i+1), 10)]; ok {

			_, hasContainer := instance.GetContainer(true)

			if !hasContainer {
				continue InstanceScaleReturn
			}

			node.log.Message("Node Scaling down. Starting instance :"+instance.Name)
			instance.Stop(false,10)

			count++
			if count >= number {
				return count
			}
		}
// operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "SCALE NODE TEST INSTANCE:", *target.node.InstanceMap[strconv.FormatInt(int64(i+1), 10)])
	}
	
	return count
}