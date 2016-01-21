package libs

/**
 * @file Build client factories from YAML files
 */

import (
	"io/ioutil"
	"strings"

	"encoding/json"
	"gopkg.in/yaml.v2"

	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

const (
	COACH_CONF_CLIENTSFILE = "clients.yml"
)

// Look for factory configurations inside the project confpaths
func (clientFactories *ClientFactories) from_ClientFactoriesYaml(logger log.Log, project *conf.Project) {
	for _, yamlConfFilePath := range project.Paths.GetConfSubPaths(COACH_CONF_CLIENTSFILE) {
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Looking for YAML clients file: "+yamlConfFilePath)
		clientFactories.from_ClientFactoriesYamlFilePath(logger, project, yamlConfFilePath)
	}
}

// Try to configure factories by parsing yaml from a conf file
func (clientFactories *ClientFactories) from_ClientFactoriesYamlFilePath(logger log.Log, project *conf.Project, yamlFilePath string) bool {
	// read the config file
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Could not read a YAML file: "+err.Error())
		return false
	}

	if !clientFactories.from_ClientFactoriesYamlBytes(logger.MakeChild(yamlFilePath), project, yamlFile) {
		logger.Warning("YAML marshalling of the YAML clients file failed [" + yamlFilePath + "]: " + err.Error())
		return false
	}
	return true
}

// Try to configure factories by parsing yaml from a byte stream
func (clientFactories *ClientFactories) from_ClientFactoriesYamlBytes(logger log.Log, project *conf.Project, yamlBytes []byte) bool {
	if project != nil {
		// token replace
		tokens := &project.Tokens
		yamlBytes = []byte(tokens.TokenReplace(string(yamlBytes)))
		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Tokenized Bytes", string(yamlBytes))
	}

	var yaml_clients map[string]map[string]interface{}
	err := yaml.Unmarshal(yamlBytes, &yaml_clients)
	if err != nil {
		logger.Warning("YAML parsing error : " + err.Error())
		return false
	}
	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "YAML source:", yaml_clients)

	for name, client_struct := range yaml_clients {
		clientType := ""
		client_json, _ := json.Marshal(client_struct)
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Single client JSON:", string(client_json))

		if clientType_struct, ok := client_struct["Type"]; ok {
			clientType, _ = clientType_struct.(string)
		} else {
			clientType = name
		}

		switch strings.ToLower(clientType) {
		case "docker":
			fallthrough
		case "fsouza":

			clientFactorySettings := &FSouza_ClientFactorySettings{}
			err := json.Unmarshal(client_json, clientFactorySettings)

			if err != nil {
				logger.Warning("Factory definition failed to configure client factory :" + err.Error())
				logger.Debug(log.VERBOSITY_DEBUG, "Factory configuration json: ", string(client_json), clientFactorySettings)
				continue
			}

			factory := FSouza_ClientFactory{}
			if !factory.Init(logger.MakeChild(clientType), project, ClientFactorySettings(clientFactorySettings)) {
				logger.Error("Failed to initialize FSouza factory from client factory configuration: " + err.Error())
				continue
			}

			// Add this factory to the factory list
			logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Client Factory Created [Client_DockerFSouzaFactory]", factory)
			clientFactories.AddClientFactory(clientType, ClientFactory(&factory))

		case "":
			logger.Warning("Client registration failure, client has a bad value for 'Type'")
		default:
			logger.Warning("Client registration failure, client has an unknown value for 'Type' :" + clientType)
		}

	}

	return true
}
