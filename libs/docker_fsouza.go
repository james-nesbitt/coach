package libs

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strings"

	docker "github.com/fsouza/go-dockerclient"

	"github.com/james-nesbitt/coach-tools/log"
)

/**
 * Client Factory
 */

// Client Factory settings
type Client_DockerFSouzaFactorySettings struct {
	Host     string `json:"Host,omitempty" yaml:"Host,omitempty"`
	CertPath string `json:"CertPath,omitempty" yaml:"CertPath,omitempty"`
}

type Client_DockerFSouzaFactory struct {
	Client_DockerFSouzaFactorySettings
	log    log.Log
	client *docker.Client
}

func (clientFactory *Client_DockerFSouzaFactory) MeetsRequirements(requirements FactoryMatchRequirements) bool {
	return true
}
func (clientFactory *Client_DockerFSouzaFactory) MakeClient(logger log.Log, configJson string) *Client_DockerFSouza {

	// if we haven't made an actual fsouza docker client, then do it now.
	if clientFactory.client == nil {
		if client, err := clientFactory.makeFsouzaClient(logger.MakeChild("fsouza")); err == nil {
			clientFactory.client = client
		} else {
			logger.Error("Failed to create actual FSouza Docker client from client factory configuration: " + err.Error())
			// @NOTE it is a bit late to fail, as we are already creating clients at this point?
		}
	}

	var settings Client_DockerFSouzaSettings
	err := json.Unmarshal([]byte(configJson), &settings)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return &Client_DockerFSouza{
		Client_DockerFSouzaSettings: settings,
		client: clientFactory.client,
	}
}
func (clientFactory *Client_DockerFSouzaFactory) MakeNodeClient(logger log.Log, configJson string) NodeClient {
	return NodeClient(clientFactory.MakeClient(logger, configJson))
}
func (clientFactory *Client_DockerFSouzaFactory) MakeInstanceClient(logger log.Log, configJson string) InstanceClient {
	return InstanceClient(clientFactory.MakeClient(logger, configJson))
}

func (clientFactory *Client_DockerFSouzaFactory) Match(requirements FactoryMatchRequirements) bool {
	return true
}

/**
 * makeFsouzaClient is a factory where we make an actualy fsouz client, which can
 * be used for as many NodeClient and InstanceClients as we need.
 * Tyipcall the Client_DockerFSouzaFactory makes a single client, and then uses it
 * repeatedly when creating Client_DockerFSouza objects.
 */
func (clientFactory *Client_DockerFSouzaFactory) makeFsouzaClient(logger log.Log) (*docker.Client, error) {
	var client *docker.Client
	var err error
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Docker client conf: ", clientFactory)

	if strings.HasPrefix(clientFactory.Host, "tcp://") {

		if _, err := os.Stat(clientFactory.CertPath); err == nil {
			// TCP DOCKER CLIENT WITH CERTS
			client, err = docker.NewTLSClient(
				clientFactory.Host,
				path.Join(clientFactory.CertPath, "cert.pem"),
				path.Join(clientFactory.CertPath, "key.pem"),
				path.Join(clientFactory.CertPath, "ca.pem"),
			)
		} else {
			// TCP DOCKER CLIENT WITHOUT CERTS
			client, err = docker.NewClient(clientFactory.Host)
		}

	} else if strings.HasPrefix(clientFactory.Host, "unix://") {
		// TCP DOCKER CLIENT WITHOUT CERTS
		client, err = docker.NewClient(clientFactory.Host)
	} else {
		err = errors.New("Unknown client host :" + clientFactory.Host)
	}

	if err == nil {
		logger.Debug(log.VERBOSITY_DEBUG_WOAH, "FSouza Docker client created:", client)
	}
	return client, err
}

/**
 * Client objects
 */

// Client Instance settings
type Client_DockerFSouzaSettings struct {
	BuildPath string            `json:"Build,omitempty" yaml:"Build,omitempty"`
	Config    docker.Config     `json:"Config,omitempty" yaml:"Config,omitempty"`
	Host      docker.HostConfig `json:"Host,omitempty" yaml:"Host,omitempty"`
}

type Client_DockerFSouza struct {
	Client_DockerFSouzaSettings

	client *docker.Client
}

/**
 * Client Settings interface: Relational Methods
 */

func (client *Client_DockerFSouzaSettings) FindRelatedNodes(nodes Nodes) {

}

func (client *Client_DockerFSouzaSettings) MakeInstanceClient(instance Instance) {

}

func (client *Client_DockerFSouza) TakeNodeSettings(settings NodeClientSettings) {

}
func (client *Client_DockerFSouza) TakeInstanceSettings(settings InstanceClientSettings) {

}

/**
 * NodeClient meta-methods
 */

func (client *Client_DockerFSouza) HasImage() bool {
	return false
}

func (client *Client_DockerFSouza) HasContainer() bool {
	return false
}

/**
 * InstanceClient meta-methods
 */

func (client *Client_DockerFSouza) IsRunning() bool {
	return false
}

/**
 * NodeClient interface: Operation Methods
 */

func (client *Client_DockerFSouza) Attach(force bool) bool {
	return false
}

func (client *Client_DockerFSouza) Build() bool {
	return false
}

func (client *Client_DockerFSouza) Destroy() bool {
	return false
}

func (client *Client_DockerFSouza) Pull() bool {
	return false
}

func (client *Client_DockerFSouza) Create() bool {
	return false
}

func (client *Client_DockerFSouza) Remove(force bool) bool {
	return false
}

func (client *Client_DockerFSouza) Start() bool {
	return false
}

func (client *Client_DockerFSouza) Stop(force bool) bool {
	return false
}

func (client *Client_DockerFSouza) Pause() bool {
	return false
}

func (client *Client_DockerFSouza) Unpause() bool {
	return false
}

func (client *Client_DockerFSouza) Info() bool {
	return false
}

func (client *Client_DockerFSouza) Commit() bool {
	return false
}

func (client *Client_DockerFSouza) Run(cmdOverride []string) bool {
	return false
}
