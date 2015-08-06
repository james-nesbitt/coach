package main

import (
	"strings"
)

/**
 * Match any Config and HostConfig identifiers to any nodes in the nodes array
 * and add them to the Node's map of dependencies
 */
func (node *Node) SetDependencies(nodes Nodes) {

	dependencies := []string{}

	if node.HostConfig.Links!=nil {
		for _, name := range node.HostConfig.Links {
			dependencies = append(dependencies, DependencyBase(name))
		}
	}
	if node.HostConfig.VolumesFrom!=nil {
		for _, name := range node.HostConfig.VolumesFrom {
			dependencies = append(dependencies, DependencyBase(name))
		}
	}

	for _, dependency := range dependencies {
		if _, ok := nodes.Map[dependency]; ok {
			node.Dependencies[dependency] = nodes.Map[dependency]
		}
	}

}

/**
 * Interpret the instance dependency format, for a node instance identifier
 *
 * This format interpreter allows nodes to use node instance identifiers in their
 * Docker settings for values such as --links, and --volumes-from.  the format
 * allows a simple synax from one node, that defines a particular target instance.
 *
 * This gets used by the node processors various times to convert syntax into
 * container name
 *
 * @param targets []string a slice of strings of the either syntax "instance" or "instance:destination", as used
 *    in the docker client VolumesFrom, Links etc fields
 *
 * @note subsitutions are made if the string target can be matched to a node in the Dependencies map,
 *   with the following syntax allowed to target specific instances:
 *     {node}                   => fallback instance of that node will be matched
 *     {node}->{instance}       => {instance} or the fallback of that node will be matched
 *     {node}->all              => all instances of that node will be matched
 *     {node}->random           => a random instance of that node will be matched (fallback is ignored)
 *
 * @returns Container Name as a string, and success boolean
 */
func (node *Node) DependencyInstanceMatches(targets []string, fallbackInstance string) []string {
	newTargets := []string{}
	var targetInstances []*Instance

	for _, target := range targets {
		targetSplit := strings.SplitN(target, ":", 2)
		targetNodeName := targetSplit[0]
		targetInstanceName := ""

		if strings.Contains(targetNodeName, "->") {
			subSplit := strings.SplitN(targetNodeName, "->", 2)
			targetNodeName = subSplit[0]
			targetInstanceName = subSplit[1]
		}

		if targetNode, ok := node.Dependencies[targetNodeName]; ok {
			targetInstances = []*Instance{}

			if targetInstanceName=="%ALL" {
				// in this case we have a meta request to link to all active instances of a target
				targetInstances = targetNode.GetInstances(true)
			} else if targetInstanceName=="%RANDOM" {
				if randomTarget := targetNode.GetRandomInstance(true); randomTarget!=nil {
					targetInstances = []*Instance{ randomTarget }
				}
			} else {
				targetInstances = targetNode.FilterInstances([]string{targetInstanceName, fallbackInstance}, true)
			}

			for _, target := range targetInstances {
				newTarget := target.GetContainerName()
				if len(targetSplit)>1 {
					if strings.Contains(targetSplit[1], "%TARGET") {
						newTarget += ":"+strings.Replace(targetSplit[1], "%TARGET", target.Name, -1)
					} else {
						newTarget += ":"+target.Name+"_"+targetSplit[1]
					}
				}
				newTargets = append(newTargets, newTarget)
			}
		}
	}

	return newTargets
}


func DependencyBase(dependency string) string {
	if strings.Contains(dependency, "->") {
		return strings.SplitN(dependency, "->", 2)[0]
	} else if strings.Contains(dependency, ":") {
		return strings.SplitN(dependency, ":", 2)[0]
	} else {
		return dependency
	}
}
