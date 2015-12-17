package client

import (
	"github.com/james-nesbitt/coach-tools/conf"
	"github.com/james-nesbitt/coach-tools/log"
	// "github.com/james-nesbitt/coach-tools/node"
)

// Client constructor
func GetClient(clientLog log.CoachLog, conf *conf.Conf) *Client {
  switch conf.GetSettings['client'] {
  	default: 
  	  // by default use the FSouza based docker client
  	  return &Client(&getClient_DockerFsouza(clientLog, conf))
  }
}

// Client for container manager connections for coach nodes
type Client interface {
	// Operation handlers
	attachInstance(instance *Instance) bool
	buildNode(node *Node) bool
	commitInstance(instance *Instance) bool
	createInstance(instance *Instance) bool
	destroyNode(node *Node) bool
	infoNode(node *Node) bool
	pauseInstance(instance *Instance) bool
	pullNode(node *Node) bool
	removeInstance(instance *Instance) bool
	unpauseInstance(instance *Instance) bool
	startInstance(instance *Instance) bool
	stopInstance(instance *Instance) bool
}
