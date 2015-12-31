package libs

import (
	"github.com/james-nesbitt/coach-tools/conf"
	"github.com/james-nesbitt/coach-tools/log"
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
		factories.from_Default(logger.MakeChild("default"))
	}

	return factories
}

/**
 * This is a fallback client builder, which builds the default
 * coach client.  The default coach client is currently the
 * FSouza Docker client, configured to use ENV settings, or
 * a local socket.
 */
func (clientFactories *ClientFactories) from_Default(logger log.Log) {

}

// ClientFactories An ordered collection of NodeClient/InstanceCLient factories
type ClientFactories struct {
	log log.Log
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

	Init(logger log.Log, settings ClientFactorySettings) bool

	MakeClient(logger log.Log, settings ClientSettings) (Client, bool)
}

type ClientFactorySettings interface {
	Settings() interface{}
}
