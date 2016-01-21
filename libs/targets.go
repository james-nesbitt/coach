package libs

import (
	"strings"

	"github.com/twmb/algoimpl/go/graph"

	"github.com/james-nesbitt/coach-tools/log"
)

// Build a targets object for a nodes list, from a list of string identifiers
func (nodes *Nodes) Targets(logger log.Log, identifiers []string) *Targets {
	targets := &Targets{log: logger, targetMap: map[string]*Target{}, targetOrder: []string{}}
	targets.fromNodes(identifiers, *nodes)
	targets.Sort()
	return targets
}

// A set of node targets
type Targets struct {
	log         log.Log
	targetMap   map[string]*Target
	targetOrder []string
}

func (targets *Targets) Target(id string) (target *Target, ok bool) {
	target, ok = targets.targetMap[id]
	return
}
func (targets *Targets) TargetOrder() []string {
	return targets.targetOrder
}

// Build up a targets list by interpreting string identifiers as a set of nodes targets
func (targets *Targets) fromNodes(identifiers []string, nodes Nodes) {
	targets.log.Debug(log.VERBOSITY_DEBUG, "Adding targets from nodes", identifiers)

	for _, identifier := range identifiers {
		prefix := identifier[0:1]
		switch prefix {
		case "$": // note that it is impossible to pass these in on the command line
			identifier = string(identifier[1:])
			switch identifier {
			case "all":
				for _, name := range nodes.NodeNames() {
					if node, ok := nodes.Node(name); ok {
						_, instances := targets.targetStringSeparate(identifier)
						targets.addNodeTarget(name, node, instances)
					}
				}
			}

		case "%":
			identifier = string(identifier[1:])
			for name, node := range nodes.NodesMap {
				if node.Type() == identifier {
					if node, ok := nodes.Node(name); ok {
						_, instances := targets.targetStringSeparate(identifier)
						targets.addNodeTarget(name, node, instances)
					}
				}
			}

		case "@":
			identifier = string(identifier[1:])
			fallthrough
		default:
			name, instances := targets.targetStringSeparate(identifier)
			if node, ok := nodes.Node(name); ok {
				targets.addNodeTarget(name, node, instances)
			} else {
				targets.log.Error("Unknown identifier passed: " + prefix + identifier)
			}

		}
	}
}

// translate a single node.instance1.instance2 string a node name, and a slice of instance filters
func (targets *Targets) targetStringSeparate(identifier string) (node string, instances []string) {
	separated := strings.Split(identifier, ":")
	if len(separated) > 1 {
		return separated[0], separated[1:]
	} else {
		return separated[0], []string{}
	}
}

// Register a node (with instances filters) as a named target
func (targets *Targets) addNodeTarget(name string, node Node, instanceFilters []string) {
	if _, exists := targets.targetMap[name]; !exists {
		// this is a new target node, so add all of the instances
		instances, _ := node.Instances().FilterableInstances()
		target := Target{name: name, node: node, instances: instances}

		targets.targetMap[name] = &target
		targets.targetOrder = append(targets.targetOrder, name)

		targets.log.Debug(log.VERBOSITY_DEBUG_LOTS, "Added target", name, target)
	}

	if instances, ok := targets.targetMap[name].Instances(); ok {
		// The target node has already been added, so we just need to confirm instances
		if len(instanceFilters) > 0 {
			instances.AddFilters(instanceFilters...)
		}
	}
}

// Sort the node targets based on dependencies (using a graph-sort)
func (targets *Targets) Sort() bool {
	logger := targets.log.MakeChild("sort")
	logger.Debug(log.VERBOSITY_DEBUG, "Starting target sort:", targets, targets.TargetOrder())

	g := graph.New(graph.Directed)
	graphNodes := make(map[string]graph.Node, 0)
	targetMap := map[string]*Target{}

	// use targets to create graph nodes
	for _, name := range targets.TargetOrder() {
		targetMap[name], _ = targets.Target(name)
		graphNodes[name] = g.MakeNode()
		*graphNodes[name].Value = name
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "ADDING TARGET:"+name, nil)
	}
	// use dependencies to set graph edges (can't be done in the above loop)
	for _, outerName := range targets.TargetOrder() {
		for _, innerName := range targets.TargetOrder() {
			outerTarget, _ := targets.Target(outerName)
			outerNode, _ := outerTarget.Node()

			if outerNode.DependsOn(innerName) {
				g.MakeEdge(graphNodes[innerName], graphNodes[outerName])
				logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "ADDING EDGE:"+innerName+"<-"+outerName, nil)
			}
		}
	}

	sorted := []string{}
	for _, graphNode := range g.TopologicalSort() {
		value := *graphNode.Value
		if name, ok := value.(string); ok {
			sorted = append(sorted, name)
		}
	}

	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "sort SORTED:", sorted)
	targets.targetOrder = sorted
	return true
}

// A single node target
type Target struct {
	name      string
	node      Node
	instances FilterableInstances
}

func (target *Target) Name() string {
	return target.name
}
func (target *Target) Node() (Node, bool) {
	return target.node, target.node != nil
}
func (target *Target) Instances() (FilterableInstances, bool) {
	return target.instances, target.instances != nil
}
