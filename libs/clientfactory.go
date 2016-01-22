package libs

import (
	"os"

	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

// Generate a factory set from a project
func MakeClientFactories(logger log.Log, project *conf.Project) *ClientFactories {
	factories := &ClientFactories{
		log: logger,
		orderedClientFactories: []ClientFactory{},
	}

	/**
	 * Build Factories from YAML if possible
	 */
	factories.from_ClientFactoriesYaml(logger.MakeChild("fromyaml"), project)

	/**
	 * If no factory is set up then assume that we
	 * should use the docker fsouza library with a
	 * local docker implementation
	 */
	if !factories.HasFactories() {
		logger.Debug(log.VERBOSITY_DEBUG, "No defined clients, retreiving default client")
		factories.from_Default(logger.MakeChild("default"), project)
	}

	return factories
}

/**
 * This is a fallback client builder, which builds the default
 * coach client.  The default coach client is currently the
 * FSouza Docker client, configured to use ENV settings, or
 * a local socket.
 */
func (clientFactories *ClientFactories) from_Default(logger log.Log, project *conf.Project) {
	clientFactorySettings := &FSouza_ClientFactorySettings{}
	clientType := "fsouza"

	if DockerHost := os.Getenv("DOCKER_HOST"); DockerHost == "" {
		logger.Debug(log.VERBOSITY_DEBUG, "No local environment DOCKER settings found, assuming a locally running docker client will be found.")
		clientFactorySettings.Host = "unix:///var/run/docker.sock"
	} else {
		clientFactorySettings.Host = DockerHost
	}

	// if we have no cert path, and we are going to use a TCP socket, test for a default cert path.
	if DockerCertPath := os.Getenv("DOCKER_CERT_PATH"); DockerCertPath != "" {
		clientFactorySettings.CertPath = DockerCertPath
	}

	factory := FSouza_ClientFactory{}
	if !factory.Init(logger, project, ClientFactorySettings(clientFactorySettings)) {
		logger.Error("Failed to initialize FSouza factory from client factory configuration")
	}

	// Add this factory to the factory list
	logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Client Factory Created [Client_DockerFSouzaFactory]", factory)
	clientFactories.AddClientFactory(clientType, ClientFactory(&factory))
}

// ClientFactories An ordered collection of NodeClient/InstanceCLient factories
type ClientFactories struct {
	log                    log.Log
	orderedClientFactories []ClientFactory
}

func (clientFactories *ClientFactories) AddClientFactory(ID string, client ClientFactory) {
	clientFactories.orderedClientFactories = append(clientFactories.orderedClientFactories, client)
}

func (clientFactories *ClientFactories) MatchClientFactory(requirements FactoryMatchRequirements) (ClientFactory, bool) {
	for _, clientFactory := range clientFactories.orderedClientFactories {
		if clientFactory.Match(requirements) {
			clientFactories.log.Debug(log.VERBOSITY_DEBUG, "Matched client factory: "+clientFactory.Id(), nil)
			return clientFactory, true
		}
	}

	clientFactories.log.Error("Failed to match client factory")
	clientFactories.log.Debug(log.VERBOSITY_DEBUG, "Requirements", requirements)
	return nil, false
}

// HasFactories is an empty test for the factories set
func (clientFactories *ClientFactories) HasFactories() bool {
	return len(clientFactories.orderedClientFactories) > 0
}

type FactoryMatchRequirements struct {
	Type  string // Type of containerization system such as docker / rkt
	Class string // Specific Client class
	ID    string // A specific Client ID (which would have been included in requirements)
}

type ClientFactory interface {
	Id() string
	Match(requirements FactoryMatchRequirements) bool

	Init(logger log.Log, project *conf.Project, settings ClientFactorySettings) bool

	MakeClient(logger log.Log, settings ClientSettings) (Client, bool)
}

type ClientFactorySettings interface {
	Settings() interface{}
}
