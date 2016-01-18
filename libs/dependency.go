package libs

import (
	"math/rand"
	"strings"
)

type Dependencies map[string]Dependency

func (dependencies Dependencies) Dependency(id string) (dependency Dependency, ok bool) {
	dependency, ok = dependencies[id]
	return
}
func (dependencies Dependencies) SetDependency(id string, dependency Dependency) {
	dependencies[id] = dependency
}
func (dependencies Dependencies) DependencyIdTranform(id string) ([]string, bool) {
	upper, lower := dependencies.SplitId(id)

	if dependency, ok := dependencies[upper]; ok {
		return dependency.Match(lower)
	} else {
		return []string{}, false
	}
}
func (dependencies *Dependencies) SplitId(id string) (upper string, lower string) {
	upper = id
	lower = ""

	if strings.Contains(id, "->") {
		splitId := strings.Split(id, "->")
		switch len(splitId) {
		case 2:
			lower = splitId[1]
			fallthrough
		case 1:
			upper = splitId[0]
		}
	}
	return
}

type Dependency interface {
	Match(id string) ([]string, bool)
}

type NodeDependency struct {
	Node
}

func (dependency *NodeDependency) Match(id string) ([]string, bool) {
	instances := dependency.Instances()

	if id == "" {
		id = "%first"
	}

	if strings.HasPrefix(id, "%") {
		switch strings.ToLower(id[1:]) {
		case "all":
			matches := []string{}
			for _, instanceID := range instances.InstancesOrder() {
				if match, ok := instances.Instance(instanceID); ok {
					matches = append(matches, match.MachineName())
				}
			}
			return matches, true

		case "random":
			if instancesOrder := instances.InstancesOrder(); len(instancesOrder)>0 {
				instanceID := instancesOrder[rand.Intn(len(instancesOrder))]
				instance, _ := instances.Instance(instanceID)
				return []string{instance.MachineName()}, true
			}
		case "first":
			if instancesOrder := instances.InstancesOrder(); len(instancesOrder)>0 {
				instanceID := instancesOrder[0]
				instance, _ := instances.Instance(instanceID)
				return []string{instance.MachineName()}, true
			}
		}

	} else {
		if match, ok := instances.Instance(id); ok {
			return []string{match.MachineName()}, true
		}
	}

	return []string{}, false
}
