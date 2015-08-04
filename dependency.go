package main

import (
	"strings"
)

func (node *Node) SetDependencies(nodes Nodes) {

	dependencies := []string{}

	if node.HostConfig.Links!=nil {
		for _, name := range node.HostConfig.Links {
			dependencies = append(dependencies, strings.SplitN(name, ":", 2)[0])
		}
	}
	if node.HostConfig.VolumesFrom!=nil {
		for _, name := range node.HostConfig.VolumesFrom {
			dependencies = append(dependencies, strings.SplitN(name, ":", 2)[0])
		}
	}

	for _, dependency := range dependencies {
		if _, ok := nodes.Map[dependency]; ok {
			node.Dependencies[dependency] = nodes.Map[dependency]
		}
	}

}
