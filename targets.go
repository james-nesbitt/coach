package main

import (
	"strings"

	"github.com/twmb/algoimpl/go/graph"
)

type Target struct {
	node *Node
	instances []*Instance
}

func addNodeToTargets(targets []*Target, node *Node, instances []*Instance) []*Target {

	for _, eachTarget := range targets {
		if eachTarget.node.Name==node.Name {
			// The target node has already been added, so we just need to confirm instances

			for _, addInstance := range instances {
				found := false
				for _, existingInstance := range eachTarget.instances {
					if addInstance.Name==existingInstance.Name {
						found = true
						break
					}
				}
				if !found {
					eachTarget.instances = append(eachTarget.instances, addInstance)
				}
			}

			return targets
		}
	}

	// this is a new target node, so add all of the instances
	target := Target{node:node, instances:instances}
	targets = append(targets, &target)
	return targets
}

func targetStringSeparate(target string) (node string, instances []string) {
	separated := strings.Split(target, ".")
	if len(separated)>1 {
		return separated[0], separated[1:]
	} else {
		return separated[0], []string{}
	}
}

func (nodes Nodes) GetTargets(targetNames []string, onlyActive bool) []*Target {
	log := nodes.log.ChildLog("TARGETS")

	targets := []*Target{}				// ordered list of targets

	log.DebugObject( LOG_SEVERITY_DEBUG, "TARGETS", targetNames)

	if len(targetNames)==0 {
		targetNames = append(targetNames, "$all")
	}

	for _, target := range targetNames {
		prefix := target[0:1]
		switch prefix {
			case "$": // note that it is impossible to pass these in on the command line
				target = string(target[1:])
				switch target {
					case "all":
						for name, _ := range nodes.Map {
							if node, ok := nodes.GetNode(name); ok {
								_, instances := targetStringSeparate(target)
								if len(instances)==0 {
									targets = addNodeToTargets(targets, node, node.GetInstances(onlyActive))
								} else {
									targets = addNodeToTargets(targets, node, node.FilterInstances(instances, onlyActive))
								}
							}
						}
				}

			case "%":
				target = string(target[1:])
				for name, node := range nodes.Map {
					if node.NodeType==target {
						if node, ok := nodes.GetNode(name); ok {
							_, instances := targetStringSeparate(target)
							if len(instances)==0 {
								targets = addNodeToTargets(targets, node, node.GetInstances(onlyActive))
							} else {
								targets = addNodeToTargets(targets, node, node.FilterInstances(instances, onlyActive))
							}
						}
					}
				}

			case "@":
				target = string(target[1:])
				fallthrough
			default:
				name, instances := targetStringSeparate(target)
				if node, ok := nodes.GetNode(name); ok {
					if len(instances)==0 {
						targets = addNodeToTargets(targets, node, node.GetInstances(onlyActive))
					} else {
						targets = addNodeToTargets(targets, node, node.FilterInstances(instances, onlyActive))
					}
				} else {
					nodes.log.Error("Unknown target passed: "+target)
				}

		}
	}

	log.DebugObject(LOG_SEVERITY_DEBUG_STAAAP, "Get Targets:", targets)

	sorted := nodes.SortTargets(targets)

	log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Get Sorted:", sorted)

	return sorted
}
func (nodes *Nodes) SortTargets(targets []*Target) []*Target {
	log := nodes.log.ChildLog("SORT")

	g := graph.New(graph.Directed)
	graphNodes := make(map[string]graph.Node, 0)
	targetMap := map[string]*Target{}

	// use targets to create graph nodes
	for _, target := range targets {
		name := target.node.Name
		targetMap[name] = target
		graphNodes[name] = g.MakeNode()
		*graphNodes[name].Value = name
		log.Debug(LOG_SEVERITY_DEBUG_STAAAP, "ADDING TARGET:"+name)
	}
	// use dependencies to set graph edges
	for _, target := range targets {
		if target.node.Dependencies!=nil {
			for dependency, _ := range target.node.Dependencies {
				if _, ok := graphNodes[dependency]; ok {
					g.MakeEdge(graphNodes[dependency], graphNodes[target.node.Name])
					log.Debug(LOG_SEVERITY_DEBUG_STAAAP, "ADDING EDGE:"+dependency+"<-"+target.node.Name)
				}
			}
		}
	}

	log.Debug( LOG_SEVERITY_DEBUG_STAAAP, "GETTING SORTED ITEMS")
	sorted := []*Target{}
	for _,graphNode := range g.TopologicalSort() {
		value := *graphNode.Value
		log.DebugObject( LOG_SEVERITY_DEBUG_STAAAP, "SORTED ITEM:", value)
		if name, ok := value.(string); ok {
			sorted = append(sorted, targetMap[name])
		}
	}

	log.DebugObject( LOG_SEVERITY_DEBUG_STAAAP, "sort SORTED:", sorted)
	return sorted
}
