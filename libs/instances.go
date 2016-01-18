package libs

import (
	"github.com/james-nesbitt/coach-tools/log"
)

const (
	INSTANCES_FILTER_DEFAULT = "$default"
	INSTANCES_FILTER_ALL     = "$all"
)

type InstancesSettings interface {
	Settings() interface{}
}

type Instances interface {
	MachineName() string

	Init(logger log.Log, machineName string, client Client, settings InstancesSettings) bool
	Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool

	Client() InstancesClient
	FilterableInstances() (FilterableInstances, bool)

	Instance(id string) (Instance, bool)
	InstancesOrder() []string
}

// A Filterable set of instances from an Instances object
type FilterableInstances interface {
	AddFilters(...string)
	UseAll()

	Instance(id string) (Instance, bool)
	InstancesOrder() []string
}

type BaseInstances struct {
	machineName string
	log         log.Log
	client      Client

	instancesMap   map[string]Instance
	instancesOrder []string
}

func (instances *BaseInstances) Init(logger log.Log, machineName string, client Client, settings InstancesSettings) bool {
	instances.machineName = machineName
	instances.log = logger
	instances.client = client
	instances.instancesMap = map[string]Instance{}
	instances.instancesOrder = []string{}
	return true
}
func (instances *BaseInstances) Prepare(logger log.Log, client Client, nodes *Nodes, node Node) bool {
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Prepare: Base Instances")
	return true
}

// Return a unique machine name for the instances
func (instances *BaseInstances) MachineName() string {
	return instances.machineName
}

// Return an InstancesClient for the instances
func (instances *BaseInstances) Client() InstancesClient {
	return instances.client.InstancesClient(instances)
}
func (instances *BaseInstances) Instance(id string) (instance Instance, ok bool) {
	instance, ok = instances.instancesMap[id]
	return
}

// Give an ordered list of string instance IDs for this instances object
func (instances *BaseInstances) InstancesOrder() []string {
	return instances.instancesOrder
}

// Give a filterable instances for this instances object
func (instances *BaseInstances) FilterableInstances() (FilterableInstances, bool) {
	filterableInstances := BaseFilterableInstances{Instances: Instances(instances), filters: []string{}}
	return FilterableInstances(&filterableInstances), true
}

// Extend the BaseInstances with a set of filters
type BaseFilterableInstances struct {
	Instances
	filters []string // Ordered string filters, usually instance ids
}

// add a filter
func (instances *BaseFilterableInstances) AddFilters(newfilters ...string) {
	// check to see if we have been to take all instances
	if len(instances.filters) > 0 && instances.filters[0] == INSTANCES_FILTER_ALL {
		return
	}

	// remove existing filters from the new list
	for _, existingFilter := range instances.filters {
		for innerIndex, newFilter := range newfilters {
			if existingFilter == newFilter {
				// remove this item, as it is already in our filter list
				newfilters = append(newfilters[:innerIndex-1], newfilters[innerIndex+1:]...)
			}
		}
	}
	// append any filters
	if len(newfilters) > 0 {
		instances.filters = append(instances.filters, newfilters...)
	}
}

// remove all filters
func (instances *BaseFilterableInstances) UseAll() {
	instances.filters = []string{INSTANCES_FILTER_ALL}
}

// Retrieve a single Instance if it matches a filtered value
func (instances *BaseFilterableInstances) Instance(id string) (Instance, bool) {
	if len(instances.filters) > 0 && instances.filters[0] == INSTANCES_FILTER_ALL {
		instance, ok := instances.Instances.Instance(id)
		return instance, ok
	} else {
		for _, filter := range instances.filters {
			if id == filter {
				instance, ok := instances.Instances.Instance(id)
				return instance, ok
			}
		}
		return nil, false
	}
}

// Give a filtered ordered list of string instance IDs for this instances object
func (instances *BaseFilterableInstances) InstancesOrder() []string {
	// get the full instance list from the parent Instances object
	instancesOrder := instances.Instances.InstancesOrder()

	// now filter the full list
	if len(instances.filters) > 0 && instances.filters[0] == INSTANCES_FILTER_ALL {
		return instancesOrder
	} else {
		filteredOrder := []string{}
		for _, filter := range instances.filters {
			for _, instance := range instancesOrder {
				if instance == filter {
					filteredOrder = append(filteredOrder, instance)
				}
			}
		}
		return filteredOrder
	}
}
