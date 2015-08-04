package main

import (
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

// Base Node

type Node struct {
	conf *Conf
	log Log
	client *docker.Client

	Name string												// key index of a node
	MachineName string								// used as build images, and for container prefixes
	NodeType string										// role that the node will play; volume, build, service, command

	BuildPath string									// path from conf.Paths[build] to the Dockerfile for this node (optional)

	InstanceType string								// what instance approach is this node: single | fixed | scaled | temporary
	InstanceMap map[string]*Instance	// an actual map of instance objects

	Dependencies map[string]*Node			// a map to other nodes which this node is possibly dependent on

	do map[string]bool 								// permissions list

	Config docker.Config							// docker client configuration for container
	HostConfig docker.HostConfig			// docker client configuration for host

	processed bool										// track if a node needs to be processed (THIS IS NOT VERY FLEXIBLE)
}

/**
 * INIT NODE: This is a processing step done when the node is
 * first created, to make sure that it has all of the needed elements
 * for operations.
 */

func (node *Node) Init(nodeType string) bool {

	if node.MachineName=="" {
		node.MachineName = node.conf.Project+"_"+node.Name
	}

	if node.InstanceMap==nil {
		node.InstanceMap = map[string]*Instance{}
	}
	if node.Dependencies==nil {
		node.Dependencies = map[string]*Node{}
	}
	if node.do==nil {
		node.do = map[string]bool{}
	}

	if (node.BuildPath=="") {
		node.do["build"] = false
		node.do["pull"] = true
	} else {
		node.do["build"] = true
		node.do["pull"] = false
	}

	node.NodeType = nodeType
	switch nodeType {
		case "build":
			node.Init_Build()
		case "volume":
			node.Init_Volume()
		case "service":
			node.Init_Service()
		case "command":
			node.Init_Command()
	}

	node.processed = false

	return true
}

/**
 * PROCESS NODE: this is a late processing of a node, done based on a
 * set of nodes, so that the node can pull all of it's dependencies,
 * without concern that some may be missing
 */
func (node *Node) Process(nodes Nodes) {
	node.log.Debug( LOG_SEVERITY_DEBUG, "Processing Node : "+node.Name)

	// scan the node for dependencies, and add them from the nodes
	node.SetDependencies(nodes)

}

func (node *Node) GetImageName() string {
	if node.Config.Image=="" {
		return strings.ToLower(node.MachineName)
	}
	return node.Config.Image
}

// check if a node should do an action (if it is permitted)
func (node *Node) Do(action string) (bool) {
	if do, ok := node.do[action]; ok {
		return do
	}
	return false
}

/**
 * Type specific Inits
 */

func (node *Node) Init_Build() {
	node.InstanceType = "none"

	node.do["create"] = false
	node.do["start"] = false
	node.do["run"] = false
	node.do["exec"] = false

	node.do["build"] = true
	node.do["pull"] = true
	node.do["commit"] = false
}

func (node *Node) Init_Volume() {
	node.do["create"] = true
	node.do["start"] = false
	node.do["run"] = false
	node.do["exec"] = false
	node.do["commit"] = true
}

func (node *Node) Init_Service() {
	node.do["create"] = true
	node.do["start"] = true
	node.do["run"] = false
	node.do["exec"] = true
	node.do["commit"] = true
}

func (node *Node) Init_Command() {
	node.InstanceType = "temporary"

	node.do["create"] = false // only temporary instances will be allowed
	node.do["start"] = false
	node.do["run"] = true
	node.do["exec"] = false
	node.do["commit"] = true
}
