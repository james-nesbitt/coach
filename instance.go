package main

import (
	"strings"
	"strconv"
	"path"

	docker "github.com/fsouza/go-dockerclient"
)

func (node *Node) ConfigureInstances_Single() bool {
	node.InstanceType = "single"
	node.AddInstance("single", true)
	return true
}
func (node *Node) ConfigureInstances_Temporary() bool {
	node.InstanceType = "temporary"
	return true
}
func (node *Node) ConfigureInstances_Fixed(instances []string) bool {
	if node.Do("start") {
		node.InstanceType = "fixed"
		node.AddInstances(instances, true)
		return true
	} else {
		return false
	}
}
func (node *Node) ConfigureInstances_Scaled(min int, max int) bool {
	if node.Do("start") {

		node.InstanceType = "scaled"
		for i:=0; i<max; i++ {
			node.AddInstance(strconv.Itoa(i+1), i<min)
		}
		return true

	} else {
		return false
	}
}

/**
 * ADDING INSTANCES : means adding a "record" of where an instance container
 * could be tracked.  The adding process does not process any instances, and
 * does not verify if any containers exist, as it needs to check dependencies
 * first, which it cannot do until all nodes have been created/added
 */

func (node *Node) AddInstance(name string, active bool) {
	if node.InstanceMap==nil {
		node.InstanceMap = map[string]*Instance{}	// an actual map of instance objects, starts empty to be filled in by the main node handler
	}

	instanceTokens := map[string]string{
		"INSTANCE":name,
	}

	instance := Instance{
		Node: node,
		Name: name,
		MachineName: node.MachineName+"_"+name,

		active: active,
	}

	/**
	 * Here we add some token-replaces copies of the node configs
	 * to the instance.  We don't process these configs until the
	 * node is processed, as it may require dependencies which have
	 * not yet been added
	 */

	config := DockerClient_Config_Copy( node.Config )
	config = ConfigTokenReplace(config, instanceTokens)

	instance.Config = config

	hostConfig := DockerClient_HostConfig_Copy( node.HostConfig )
	hostConfig = HostConfigTokenReplace(hostConfig, instanceTokens)

	instance.HostConfig = hostConfig

	instance.Init()
	node.InstanceMap[name] = &instance
}
func (node *Node) AddTemporaryInstance(name string) {
	node.AddInstance(name, true)
}
func (node *Node) AddInstances(instances []string, active bool) {
	for _, instance := range instances {
		node.AddInstance( instance, active )
	}
}


func (node *Node) GetInstances(onlyActive bool) []*Instance {
	instances := []*Instance{}
	for name, _ := range node.InstanceMap {
		instance := node.GetInstance(name)

		if onlyActive && !instance.active {
			continue
		}

		instances = append(instances, instance)
	}
	return instances
}
func (node *Node) FilterInstances(filters []string, onlyActive bool) []*Instance {
	instances := []*Instance{}
	for name, _ := range node.InstanceMap {
		instance := node.GetInstance(name)

		if onlyActive && !instance.active {
			continue
		}

		for _, filter := range filters {
			if filter==name {
				instances = append(instances, instance)
				break;
			}
		}
	}
	return instances
}
func (node *Node) GetInstance(name string) *Instance {

	// shortcut to get any random instance
	if name=="" {
		for _, instance := range node.InstanceMap {
			if instance.active {
				return instance
			}
		}
		// there are no instances
		return nil
	}

	// get a specific names instance
	if instance, ok := node.InstanceMap[name]; ok {
		if instance.Process() {
			return instance
		}
	}

	// no matching instance found
	return nil
}


// A single instance of a node

type Instance struct {
	Node *Node

	Name string
	MachineName string

	Config docker.Config
	HostConfig docker.HostConfig

	active bool											// should this instance be active for operations (or is it dormant, perhaps for scaling)
	processed bool									// has this instance run .Process()
}
func (instance *Instance) Init() bool {
	instance.processed = false
	return true
}
func (instance *Instance) Process() bool {
	if instance.processed {
		return true
	}

	instance.Config = instance.Node.instanceConfig(instance.Name)
	instance.HostConfig = instance.Node.instanceHostConfig(instance.Name)

	instance.processed = true
	return true
}
func (instance *Instance) GetContainerName() string {
	return strings.ToLower(instance.MachineName)
}

/**
 * Elements in the nodes struct are used directly as docker configuration. but are keyed
 * to the nodes map keys.  To properly use them, these mapped elements have to be changed
 * from the Nodes key, to the proper instance container or image name.
 * An example is the node.HostConfig.Volumes_from slice, or the node.HostConfig.Binds slice.
 */
func (node *Node) instanceConfig(name string) docker.Config {
	instance, _ := node.InstanceMap[name]
	config := instance.Config


	node.log.DebugObject( LOG_SEVERITY_DEBUG_LOTS, "TRANSFORMED CONFIG FOR INSTANCE ["+name+"]: ", config)
	return config
}
func (node Node) instanceHostConfig(name string) docker.HostConfig {
	instance, _ := node.InstanceMap[name]
	config := instance.HostConfig

	// Remap all binds so that they are based on the project folder
	// - any relative path is remapped to be from the project root
	// - any path starting with ~ is remapped to the user home flder
	// - any absolute paths are left as is
	//
	if config.Binds!=nil {
		var binds []string
		for index, bind := range config.Binds {
			binds = strings.SplitN(bind, ":", 2)
			if (path.IsAbs(binds[0])) {
				continue
			} else if (binds[0][0:1]=="~") { // one would think that path can handle such terminology
				binds[0] = path.Join(node.conf.Paths["userhome"], binds[0][1:])
			} else {
				binds[0] = path.Join(node.conf.Paths["project"], binds[0])
			}
			config.Binds[index] = strings.Join(binds, ":")
		}

	}
	// convert all links from Nodes index name to container name
	if config.Links!=nil {
		for index, link := range config.Links {
			links := strings.SplitN(link, ":", 2)
			if name, ok := node.GetDependencyInstanceContainerName(links[0], name, false); ok {
				links[0] = name
				config.Links[index] = strings.Join(links,":")
			}
		}
	}
	// convert all volumes from from Nodes index name to container name
	if config.VolumesFrom!=nil {
		for index, volumesFrom := range config.VolumesFrom {
			split := strings.SplitN(volumesFrom, ":", 2)
			if name, ok := node.GetDependencyInstanceContainerName(split[0], name, false); ok {
				split[0] = name
				config.VolumesFrom[index] = strings.Join(split, ":")
			}
		}
	}

	node.log.DebugObject( LOG_SEVERITY_DEBUG_LOTS, "TRANSFORMED HOSTCONFIG FOR INSTANCE ["+name+"]: ", config)
	return config
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
 * {node}, {instance} 							: get {instance} of {node}
 * {node}@{instance}, {fallback}		: get {instance} of {node}, fallback to the {fallback} of {node}
 *
 * if strict==false, then the first instance of a node is returned if no match can be made
 *
 * @returns Container Name as a string, and success boolean
 */
func (node *Node) GetDependencyInstanceContainerName(identifier string, fallback string, strict bool) (string, bool) {
	node.log.Debug( LOG_SEVERITY_DEBUG_STAAAP, "DEPENDENCY CONTAINER SEARCH ["+identifier+"]["+fallback+"] => START")
	split := strings.SplitN(identifier, "@", 2)
	var targetNodeName, targetInstance string
	if len(split)>1 {
		targetNodeName = split[0]
		targetInstance = split[1]
	} else {
		targetNodeName = split[0]
		targetInstance = fallback
	}

	var targetNode *Node
	var instance *Instance
	var ok bool

	if targetNode, ok = node.Dependencies[targetNodeName]; !ok {
		node.log.Debug( LOG_SEVERITY_DEBUG_STAAAP, "DEPENDENCY CONTAINER SEARCH ["+identifier+"] => NO NODE FOUND")
		return "", false
	}
	node.log.Debug( LOG_SEVERITY_DEBUG_STAAAP, "DEPENDENCY NODE FOUND ["+identifier+"]["+fallback+"]["+targetNode.InstanceType+"] => "+targetNode.Name)

	// look for the particular instance
	instance = targetNode.GetInstance(targetInstance)
	if instance!=nil {
		node.log.Debug( LOG_SEVERITY_DEBUG_STAAAP, "DEPENDENCY CONTAINER SEARCH ["+identifier+"] => TARGET INSTANCE FOUND : "+instance.Name+":"+instance.GetContainerName())
		return instance.GetContainerName(), true
	}

	if targetNode.InstanceType=="single" {
		instance := targetNode.GetInstance("single")
		if instance!=nil {
			node.log.Debug( LOG_SEVERITY_DEBUG_STAAAP, "DEPENDENCY CONTAINER SEARCH ["+identifier+"] => SINGLE INSTANCE FOUND : "+instance.Name+":"+instance.GetContainerName())
			return instance.GetContainerName(), true
		}
	}

	// look for the fallback instance
	if (targetInstance!=fallback && fallback!="") {
		instance := targetNode.GetInstance(fallback)
		if instance!=nil {
			node.log.Debug( LOG_SEVERITY_DEBUG_STAAAP, "DEPENDENCY CONTAINER SEARCH ["+identifier+"] => FALLBACK INSTANCE FOUND : "+instance.Name+":"+instance.GetContainerName())
			return instance.GetContainerName(), true
		}
	}

	// use a random instance if we couldn't find a match, and we are not in strict mode
	if !strict {
		instance := targetNode.GetInstance("*random*")
		if instance!=nil {
			node.log.Debug( LOG_SEVERITY_DEBUG_STAAAP, "DEPENDENCY CONTAINER SEARCH ["+identifier+"] => RANDOM INSTANCE FOUND : "+instance.Name+":"+instance.GetContainerName())
			return instance.GetContainerName(), true
		}
	}

	node.log.Debug( LOG_SEVERITY_DEBUG_STAAAP, "DEPENDENCY CONTAINER SEARCH ["+identifier+"]["+fallback+"] => NOT FOUND")
	return "", false
}


