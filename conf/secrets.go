package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/james-nesbitt/coach-tools/log"
)

const (
	COACH_CONF_SECRETS_SUBPATH = "secrets/secrets.yml"
)

// Load tokens from any conf subpath
func (project *Project) from_SecretsYaml(logger log.Log) {
	for _, yamlSecretsFilePath := range project.Paths.GetConfSubPaths(COACH_CONF_SECRETS_SUBPATH) {
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Looking for YAML secrets file: "+yamlSecretsFilePath)
		project.from_SecretsYamlFilePath(logger, yamlSecretsFilePath)
	}
}

// Try to configure a project by parsing yaml from a conf file
func (project *Project) from_SecretsYamlFilePath(logger log.Log, yamlFilePath string) bool {
	// read the config file
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Could not read a YAML secrets file: "+err.Error())
		return false
	}

	if !project.from_SecretsYamlBytes(logger.MakeChild(yamlFilePath), yamlFile) {
		logger.Warning("YAML marshalling of the YAML secrets file failed [" + yamlFilePath + "]: " + err.Error())
		return false
	}
	return true
}

// Try to configure a project by parsing yaml secrets from a byte stream
func (project *Project) from_SecretsYamlBytes(logger log.Log, yamlBytes []byte) bool {
	// parse the config file contents as a ConfSource_projectyaml object
	source := new(secrets_Yaml)

	if err := yaml.Unmarshal(yamlBytes, source); err != nil {
		logger.Warning("YAML parsing error : " + err.Error())
		return false
	}
	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "YAML secrets source:", *source)

	return source.configureProject(logger, project)
}

// Secrets in Yaml format
type secrets_Yaml struct {
	Secrets map[string]string
}

func (secrets *secrets_Yaml) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&secrets.Secrets)
}

// Use the secrets yaml object to configure a project
func (secrets *secrets_Yaml) configureProject(logger log.Log, project *Project) bool {
	for key, value := range secrets.Secrets {
		project.SetToken(key, value)
	}

	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Configured project from YAML secrets")
	return true
}
