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

	// default scale value
	operation.scale = 1

	if len(flags)>0 {
		if flags[0]=="up" {
			operation.scale = 1
		} else if flags[0]=="down" {
			operation.scale = -1
		}
	}

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
	operation.log.Info("running scale operation")

	if operation.scale==0 {
		operation.log.Warning("scale operation was told to scale to 0")
		return
	}

	TargetScaleReturn:
	for _, target := range operation.nodes.GetTargets(operation.targets) {

		if target.node.InstanceType!="scaled" {
			operation.log.Info("Tried to scale non-scaleable node :"+target.node.Name)
			continue TargetScaleReturn
		}

		if operation.scale>0 {
			count := target.node.ScaleUpNumber(operation.scale)
			
			if count==0 {
				target.node.log.Warning(target.node.Name+": Scale operation could not scale up any new instances of node")
			} else if count<operation.scale {
				target.node.log.Warning(target.node.Name+": Scale operation could not scale up all of the requested instances of node. "+strconv.FormatInt(int64(count+1), 10)+" started.")
			} else {
				target.node.log.Message(target.node.Name+": Scale operation scaled up "+strconv.FormatInt(int64(count), 10)+" instances")
			}

		} else {
			count := target.node.ScaleDownNumber(-operation.scale)

			if count==0 {
				target.node.log.Warning(target.node.Name+": Scale operation could not scale down any new instances of node")
			} else if count<(-operation.scale) {
				target.node.log.Warning(target.node.Name+": Scale operation could not scale down all of the requested instances of node. "+strconv.FormatInt(int64(count+1), 10)+" stopped.")
			} else {
				target.node.log.Message(target.node.Name+": Scale operation scaled down "+strconv.FormatInt(int64(count), 10)+" instances")
			}

		}

	}

	cache.refresh(false, true)
}

// scale a node up a certain number of instances
func (node *Node) ScaleUpNumber(number int) int {
	count := 0

	// find the first non-running instance, and start it
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

			node.log.Info(node.Name+": Node Scaling up. Starting instance :"+instance.Name)
			instance.Start(false)

			count++
			if count >= number {
				return count
			}
		}
	}
	
	return count
}
// scale a node down a certain number of instances
func (node *Node) ScaleDownNumber(number int) int {

	count := 0

	InstanceScaleReturn:
	for i := len(node.InstanceMap); i>=0; i-- {
		if instance, ok := node.InstanceMap[strconv.FormatInt(int64(i+1), 10)]; ok {

			_, hasContainer := instance.GetContainer(true)

			if !hasContainer {
				continue InstanceScaleReturn
			}

			node.log.Info(node.Name+": Node Scaling down. Stopping instance :"+instance.Name)
			instance.Stop(false,10)

			count++
			if count >= number {
				return count
			}
		}
	}
	
	return count
}
