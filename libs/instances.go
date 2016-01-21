package libs

import (
	"github.com/james-nesbitt/coach/log"
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
	Init(coreInstances Instances, defaultIds []string) bool

	UseAll()
	IsAll() bool

	UseDefault()
	IsDefault() bool

	AddFilters(...string)
	IsFiltered() bool

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

	useAll bool

	defaultFilters []string
	useDefault     bool

	filters []string // Ordered string filters, usually instance ids
}

// Initialize filterable instances
func (instances *BaseFilterableInstances) Init(coreInstances Instances, defaultFilters []string) bool {
	instances.Instances = coreInstances
	instances.defaultFilters = defaultFilters
	instances.filters = []string{}

	instances.useDefault = true
	instances.useAll = false

	return true
}

// add a filter
func (instances *BaseFilterableInstances) AddFilters(newfilters ...string) {
	instances.useDefault = false
	instances.useAll = false

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
func (instances *BaseFilterableInstances) IsFiltered() bool {
	return len(instances.filters) > 0
}

// remove all filters
func (instances *BaseFilterableInstances) UseAll() {
	instances.filters = []string{}
	instances.useAll = true
	instances.useDefault = false
}
func (instances *BaseFilterableInstances) IsAll() bool {
	return instances.useAll
}

// use only default filters
func (instances *BaseFilterableInstances) UseDefault() {
	instances.filters = []string{}
	instances.useAll = false
	instances.useDefault = true
}
func (instances *BaseFilterableInstances) IsDefault() bool {
	return instances.useDefault
}

// Retrieve a single Instance if it matches a filtered value
func (instances *BaseFilterableInstances) Instance(id string) (Instance, bool) {
	if instances.useAll {
		instance, ok := instances.Instances.Instance(id)
		return instance, ok
	} else {
		var filters []string
		if instances.useDefault {
			filters = instances.defaultFilters
		} else {
			filters = instances.filters
		}
		for _, filter := range filters {
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
	if instances.useAll {
		return instances.Instances.InstancesOrder()
	} else {
		var filters []string
		if instances.useDefault {
			filters = instances.defaultFilters
		} else {
			filters = instances.filters
		}

		// get the full instance list from the parent Instances object
		instancesOrder := instances.Instances.InstancesOrder()
		filteredOrder := []string{}
		for _, filter := range filters {
			for _, instance := range instancesOrder {
				if instance == filter {
					filteredOrder = append(filteredOrder, instance)
				}
			}
		}
		return filteredOrder
	}
}
