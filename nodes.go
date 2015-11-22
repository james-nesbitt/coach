package main

import (
	docker "github.com/fsouza/go-dockerclient"
)

const (
  NODE_BUILD = iota
  NODE_VOLUME
  NODE_SERVICE
  NODE_COMMAND
)

/**
 * Nodes factory method
 */
func getNodes(log Log, conf *Conf, client *docker.Client, targets []string) Nodes {
	nodes := Nodes{
		conf: conf,
		log: log,

		Map: map[string]*Node{},
		client: client,
	}

	nodes.from_Yaml( log.ChildLog("YAML"), conf )
	nodes.from_DockerCompose(log.ChildLog("SECRETS"), conf )

	return nodes;
}


// A collection of nodes

type Nodes struct {
	conf *Conf
	log Log

	Map map[string]*Node
	client *docker.Client

	processed bool
}
func (nodes *Nodes) Process() bool {
	if nodes.processed {
		return true
	}
	for _, node := range nodes.Map {
		node.Process(*nodes)
	}

	nodes.processed = true
	return true
}

func (nodes *Nodes) AddNode(key string, node *Node) error {
	if node.client==nil {
		node.client = nodes.client
	}
	if node.conf==nil {
		node.conf = nodes.conf
	}
// 	if &node.log==nil {
		node.log = nodes.log.ChildLog("NODE:"+key)
// 	}

	nodes.Map[key] = node
	return nil
}
func (nodes *Nodes) GetNode(name string) (*Node, bool) {
	if !nodes.processed {
		nodes.Process()
	}
	if node, ok := nodes.Map[name]; ok {
		return node, true
	}

	return nil, false
}
