package libs

import (
	"errors"
	"os"
	"path"
	"strings"
"fmt"
	docker "github.com/fsouza/go-dockerclient"

	"github.com/james-nesbitt/coach-tools/log"
)

/**
 * FSOUZA: Settings
 */

// Client Factory settings
type Client_DockerFSouzaFactorySettings struct {
	Host     string `json:"Host,omitempty" yaml:"Host,omitempty"`
	CertPath string `json:"CertPath,omitempty" yaml:"CertPath,omitempty"`
}
func (settings *Client_DockerFSouzaFactorySettings) Settings() interface{} {
	fmt.Println("RETURNING Client_DockerFSouzaFactorySettings")
	return settings
}

// Client Instance settings
type Client_DockerFSouzaSettings struct {
	BuildPath string            `json:"Build,omitempty" yaml:"Build,omitempty"`
	Config    docker.Config     `json:"Config,omitempty" yaml:"Config,omitempty"`
	Host      docker.HostConfig `json:"Host,omitempty" yaml:"Host,omitempty"`
}
func (settings *Client_DockerFSouzaSettings) Settings() interface{} {
	return settings
}

/**
 * FSOUZA: ClientFactory
 */

type Client_DockerFSouzaFactory struct {
	settings *Client_DockerFSouzaFactorySettings
	log    log.Log
	client *docker.Client
}
func (clientFactory *Client_DockerFSouzaFactory) Id() string {
	return "Docker:FSouza"
}
func (clientFactory *Client_DockerFSouzaFactory) Match(requirements FactoryMatchRequirements) bool {
	clientFactory.log.Debug(log.VERBOSITY_DEBUG_STAAAP, "Match test for FSouza client factory:", requirements.Type, (requirements.Type=="docker"))
	return requirements.Type=="docker" || requirements.ID==clientFactory.Id() || requirements.Class=="Client_DockerFSouzaFactory"
}

func (clientFactory *Client_DockerFSouzaFactory) Init(logger log.Log, settings ClientFactorySettings) bool {
	clientFactory.log = logger

	// make sure that the settings that were given, where the proper "Client_DockerFSouzaFactory" type
	typedSettings := settings.Settings()
	switch asserted := typedSettings.(type) {
	case *Client_DockerFSouzaFactorySettings:
		clientFactory.settings = asserted
	default:
		logger.Error("Invalid settings type passed to Fsouza Factory")
		logger.Debug(log.VERBOSITY_DEBUG, "Settings passed:", asserted)
	}

	// if we haven't made an actual fsouza docker client, then do it now
	if clientFactory.client == nil {
		if client, err := clientFactory.makeBackendFsouzaClient(logger.MakeChild("fsouza")); err == nil {
			clientFactory.client = client
			return true
		} else {
			logger.Error("Failed to create actual FSouza Docker client from client factory configuration: " + err.Error())
			return false
		}
	}
	return true
}
func (factory *Client_DockerFSouzaFactory) MakeClient(logger log.Log, settings ClientSettings) (Client, bool) {
	client := &Client_DockerFSouza{}
	return Client(client), client.Init(logger, settings)
}

/**
 * makeBackendFsouzaClient is a factory where we make an actualy fsouza client, which can
 * be used for as many NodeClient and InstanceClients as we need.
 * Tyipcall the Client_DockerFSouzaFactory makes a single client, and then uses it
 * repeatedly when creating Client_DockerFSouza objects.
 */
func (clientFactory *Client_DockerFSouzaFactory) makeBackendFsouzaClient(logger log.Log) (*docker.Client, error) {
	var client *docker.Client
	var err error
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Docker client conf: ", clientFactory)

	if strings.HasPrefix(clientFactory.settings.Host, "tcp://") {

		if _, err := os.Stat(clientFactory.settings.CertPath); err == nil {
			// TCP DOCKER CLIENT WITH CERTS
			client, err = docker.NewTLSClient(
				clientFactory.settings.Host,
				path.Join(clientFactory.settings.CertPath, "cert.pem"),
				path.Join(clientFactory.settings.CertPath, "key.pem"),
				path.Join(clientFactory.settings.CertPath, "ca.pem"),
			)
		} else {
			// TCP DOCKER CLIENT WITHOUT CERTS
			client, err = docker.NewClient(clientFactory.settings.Host)
		}

	} else if strings.HasPrefix(clientFactory.settings.Host, "unix://") {
		// TCP DOCKER CLIENT WITHOUT CERTS
		client, err = docker.NewClient(clientFactory.settings.Host)
	} else {
		err = errors.New("Unknown client host :" + clientFactory.settings.Host)
	}

	if err == nil {
		logger.Debug(log.VERBOSITY_DEBUG_WOAH, "FSouza Docker client created:", client)
	}
	return client, err
}

/**
 * Client objects
 */

/**
 * FSOUZA: Client
 */

type Client_DockerFSouza struct {
	log log.Log
	settings *Client_DockerFSouzaSettings
	client *docker.Client
}

func (client *Client_DockerFSouza) Init(logger log.Log, settings ClientSettings) bool {
	client.log = logger

	// make sure that the settings that were given, where the proper "Client_DockerFSouzaFactory" type
	Settingsd := settings.Settings()
	switch asserted := Settingsd.(type) {
	case *Client_DockerFSouzaSettings:
		client.settings = asserted
		return true
	default:
		logger.Error("Invalid settings type passed to Fsouza Client")
		return false
	}
}
func (client *Client_DockerFSouza) Prepare(logger log.Log, nodes *Nodes, node Node) bool {
	return true
}

func (client *Client_DockerFSouza) NodeClient(node Node) NodeClient {
	return NodeClient(client)
}
func (client *Client_DockerFSouza) InstanceClient(instance Instance) InstanceClient {
	return InstanceClient(client)
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
