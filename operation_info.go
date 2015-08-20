package main

import (
	"strconv"
	"strings"
	"text/tabwriter"
)

type Operation_Info struct {
	log Log

	Nodes Nodes
	Targets []string
}
func (operation *Operation_Info) Flags(flags []string) {
	for _, flag := range flags {
		switch flag {

		}
	}
}

func (operation *Operation_Info) Help() {

}

func (operation *Operation_Info) Run() {

	operation.log.Message("running info operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.Targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:INFO")
	operation.Nodes.Info(operation.Targets)
}

func (nodes *Nodes) Info(targets []string) {
	for _, target := range nodes.GetTargets(targets) {
		target.log = nodes.log.ChildLog("NODE:"+target.Name)
		target.Info()
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
			container, exists := instance.GetContainer(false)

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
			if exists {
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
