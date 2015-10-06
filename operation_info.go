package main

import (
	"strconv"
	"strings"
	"text/tabwriter"
)

type Operation_Info struct {
	log Log

	nodes Nodes
	targets []string
}
func (operation *Operation_Info) Flags(flags []string) {
	for _, flag := range flags {
		switch flag {

		}
	}
}

func (operation *Operation_Info) Help(topics []string) {
	operation.log.Note(`Operation: INFO

Coach will attempt to provide project information by investigating target images and containers.

SYNTAX:
    $/> coach {targets} info

  {targets} what target nodes the operation should process ($/> coach help targets)

`)
}

func (operation *Operation_Info) Run() {

	operation.log.Message("running info operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:INFO")
	operation.nodes.Info(operation.targets)
}

func (nodes *Nodes) Info(targets []string) {
	for _, target := range nodes.GetTargets(targets, false) {
		target.node.log = nodes.log.ChildLog("NODE:"+target.node.Name)
		target.node.Info()
	}
}

func (node *Node) Info() bool {

	node.Info_Images()

	node.Info_Instances()

	return true
}

func (node *Node) Info_Images() bool {
	images := node.GetImages()

	if len(images)==0 {
		node.log.Message("Node has no Images")
	} else {
		node.log.Message("Node Images")

		w := new(tabwriter.Writer)
		w.Init(node.log, 8, 12, 2, ' ', 0)

		row := []string{
			"",
			"ID",
			"RepoTags",
			"Created",
	//		"Size",
	//		"VirtualSize",
	//	"ParentID",
	//		"RepoDigests",
	//		"Labels",
		}
		w.Write([]byte(strings.Join(row, "\t")+"\n"))

		for index, image := range images {
			row := []string{
				strconv.FormatInt(int64(index+1), 10)+":",
				image.ID[:11],
				strings.Join(image.RepoTags, ","),
				strconv.FormatInt(image.Created, 10),
	//			strconv.FormatInt(image.Size, 10),
	//			strconv.FormatInt(image.VirtualSize, 10),
	//			image.ParentID,
	// 			strings.Join(image.RepoDigests, "\n"),
	// 			strings.Join(image.Labels, "\n"),
			}
			w.Write([]byte( strings.Join(row, "\t")+"\n"))
		}
		w.Flush()
	}

	return false
}
func (node *Node) Info_Instances() bool {

	if len(node.InstanceMap)==0 {
		node.log.Message("Node has no instances")
	} else {
		node.log.Message("Node Instances TYPE:"+node.InstanceType)

		w := new(tabwriter.Writer)
		w.Init(node.log, 8, 12, 2, ' ', 0)

		row := []string{
			"",
			"Name",
			"Container",
			"Active",
			"Status",
			"ID",
			"Created",
			"Names",
		}
		w.Write([]byte(strings.Join(row, "\t")+"\n"))

		instances := node.GetInstances(false)

		for index, instance := range instances {			
			row := []string{
				strconv.FormatInt(int64(index+1), 10),
				instance.Name,
				instance.GetContainerName(),
			}
			if instance.active {
				row = append(row, "yes")
			} else {
				row = append(row, "no")
			}

			container, found := instance.GetContainer(false)
			if found {
				row = append(row,
					container.Status,
					container.ID[:12],
					strconv.FormatInt(int64(container.Created), 10),
					strings.Join(container.Names, ", "),
				)
			} else {
				row = append(row, "n/a")
			}

			w.Write([]byte( strings.Join(row, "\t")+"\n"))
		}
		w.Flush()

	}

	return true
}

func (instance *Instance) Info() bool {

	return false
}
