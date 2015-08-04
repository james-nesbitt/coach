package main

import (
	"github.com/twmb/algoimpl/go/graph"
)

func (nodes Nodes) GetTargets(targets []string) []*Node {
	log := nodes.log.ChildLog("TARGETS")

	Targets := []*Node{}				// ordered list of targets
	added := map[string]bool{}	// track which keys have already been added, to prevent repeats

log.DebugObject( LOG_SEVERITY_DEBUG, "TARGETS", targets)
	if len(targets)==0 {
		targets = append(targets, "$all")
	}

	for _, target := range targets {
		prefix := target[0:1]
		switch prefix {
			case "$": // note that it is impossible to pass these in on the command line
				switch target[1:] {
					case "all":
						for name, _ := range nodes.Map {
							if _, ok := added[target]; ok {
								// already added this target
								continue
							}
							if node, ok := nodes.GetNode(name); ok {
								Targets = append(Targets, node)
								added[name] = true
							}
						}
				}

			case "%":
				target = string(target[1:])
				for name, node := range nodes.Map {
					if node.NodeType==target {
						if _, ok := added[target]; ok {
							// already added this target
							continue
						}
						if node, ok := nodes.GetNode(name); ok {
							Targets = append(Targets, node)
							added[name] = true
						}
					}
				}

			case "@":
				target = string(target[1:])
				fallthrough
			default:
				if _, ok := nodes.Map[target]; ok {
					if _, ok := added[target]; ok {
						// already added this target
						continue
					}
					if node, ok := nodes.GetNode(target); ok {
						Targets = append(Targets, node)
						added[target] = true
					}
				} else {
					nodes.log.Error("Unknown target passed: "+target)
				}

		}
	}

	log.DebugObject(LOG_SEVERITY_DEBUG_STAAAP, "Get Targets:", Targets)

	sorted := nodes.SortTargets(Targets)

	log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Get Sorted:", sorted)

	return sorted
}
func (nodes *Nodes) SortTargets(targets []*Node) []*Node {
	log := nodes.log.ChildLog("SORT")

	g := graph.New(graph.Directed)
	graphNodes := make(map[string]graph.Node, 0)

	// use targets to create graph nodes
	for _, target := range targets {
		graphNodes[target.Name] = g.MakeNode()
		*graphNodes[target.Name].Value = target.Name
		log.Debug(LOG_SEVERITY_DEBUG_STAAAP, "ADDING TARGET:"+target.Name)
	}
	// use dependencies to set graph edges
	for _, target := range targets {
		if target.Dependencies!=nil {
			for dependency, _ := range target.Dependencies {
				if _, ok := graphNodes[dependency]; ok {
					g.MakeEdge(graphNodes[dependency], graphNodes[target.Name])
					log.Debug(LOG_SEVERITY_DEBUG_STAAAP, "ADDING EDGE:"+dependency+"<-"+target.Name)
				}
			}
		}
	}

	log.Debug( LOG_SEVERITY_DEBUG_STAAAP, "GETTING SORTED ITEMS")
	sorted := []*Node{}
	for _,graphNode := range g.TopologicalSort() {
		value := *graphNode.Value
		log.DebugObject( LOG_SEVERITY_DEBUG_STAAAP, "SORTED ITEM:", value)
		if name, ok := value.(string); ok {
			sorted = append(sorted, nodes.Map[name])
		}
	}

	log.DebugObject( LOG_SEVERITY_DEBUG_STAAAP, "sort SORTED:", sorted)
	return sorted
}
