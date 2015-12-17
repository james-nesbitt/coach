package client

import (
	"os"
	"strings"
	"path"
	"errors"

	"encoding/json" // we use this to convert projectProjectig and HostprojectProjectig to strings for token replacement.  It could easily be replaced.

	docker "github.com/fsouza/go-dockerclient"

	"github.com/james-nesbitt/coach-tools/project"
	"github.com/james-nesbitt/coach-tools/log"	
)

func GetClient_DockerFsouza(log log.Log, project *project.projectProject) Client_DockerFsouza {


  return Client_DockerFsouza{
		client: makeFsouzaClient(log, project),

		log: clientLog,

		images: []docker.APIImages{},
		containers: []docker.APIContainers{},
  }
}

func makeFsouzaClient(log log.Log, project project.projectProject) (*docker.Client, bool) {

	var client *docker.Client
	var err error 
	log.Debug(log.VERBOSITY_DEBUG_WOAH,"Docker client project: ",project.Docker)

	if (strings.HasPrefix(project.Docker.Host, "tcp://")) {

		if _, err := os.Stat(project.Docker.CertPath); err == nil {

			// TCP DOCKER CLIENT WITH CERTS
			client, err = docker.NewTLSClient(
				project.Docker.Host,
				path.Join(project.Docker.CertPath, "cert.pem"),
				path.Join(project.Docker.CertPath, "key.pem"),
				path.Join(project.Docker.CertPath, "ca.pem"),
			)

		} else {

			// TCP DOCKER CLIENT WITHOUT CERTS
			client, err = docker.NewClient(project.Docker.Host)

		}

	} else if (strings.HasPrefix(project.Docker.Host, "unix://")) {

		if _, err := os.Stat(project.Docker.Host[7:]); err != nil {
			log.Error("Docker socket does not exist: ["+project.Docker.Host+"] "+err.Error())
		} else {
			client, err = docker.NewClient(project.Docker.Host)
		}

	} else {

		err = errors.New("Unknown client host :"+project.Docker.Host)

	}

	if err != nil {
		log.Critical(err.Error())
	}

	return client, err==nil
}
}

type Client_DockerFsouza struct {
	client *docker.Client

	log log.Log

  // cached data
	images []docker.APIImages
	containers []docker.APIContainers
}

func (client *Client_DockerFsouza) clearCache(clearImages bool, clearContainers bool) {
  
}

// Information handlers
func (client *Client_DockerFsouza) CacheClear() {
	
}

func (client *Client_DockerFsouza) getImages(node *Node) bool {

}
func (client *Client_DockerFsouza) hasImage(node *Node) bool {
	
}
func (client *Client_DockerFsouza) getContainers(instance *Instance) bool {

}
func (client *Client_DockerFsouza) hasContainer(instance *Instance) bool {

}


/**
 * OPERATIONS
 */


func (client *Client_DockerFsouza) attachInstance(instance *Instance) bool {
	
}


func (client *Client_DockerFsouza) buildNode(node *Node) bool {
	
}


func (client *Client_DockerFsouza) commitInstance(instance *Instance) bool {
	
}



func (client *Client_DockerFsouza) createInstance(instance *Instance) bool {
	
}


func (client *Client_DockerFsouza) destroyNode(node *Node) bool {
	
}


func (client *Client_DockerFsouza) infoNode(node *Node) bool {
	
}


func (client *Client_DockerFsouza) pauseInstance(instance *Instance) bool {
	
}


func (client *Client_DockerFsouza) pullNode(node *Node) bool {
	
}


func (client *Client_DockerFsouza) removeInstance(instance *Instance) bool {
	
}


func (client *Client_DockerFsouza) unpauseInstance(instance *Instance) bool {
	
}


func (client *Client_DockerFsouza) startInstance(instance *Instance) bool {
	
}


func (client *Client_DockerFsouza) stopInstance(instance *Instance) bool {
	
}
