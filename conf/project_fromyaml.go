package conf

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/james-nesbitt/coach-tools/log"
)

const (
	COACH_CONF_YAMLFILE = "conf.yml"
)

// Look for project configurations inside the project confpaths
func (project *Project) from_ConfYaml(logger log.Log) {
	for _, yamlConfFilePath := range project.Paths.GetConfSubPaths(COACH_CONF_YAMLFILE) {
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Looking for YAML conf file: "+yamlConfFilePath)
		project.from_ConfYamlFilePath(logger, yamlConfFilePath)
	}
}

// Try to configure a project by parsing yaml from a conf file
func (project *Project) from_ConfYamlFilePath(logger log.Log, yamlFilePath string) bool {
	// read the config file
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Could not read a YAML file: "+err.Error())
		return false
	}

	if !project.from_ConfYamlBytes(logger.MakeChild(yamlFilePath), yamlFile) {
		logger.Warning("YAML marshalling of the YAML conf file failed [" + yamlFilePath + "]: " + err.Error())
		return false
	}
	return true
}

// Try to configure a project by parsing yaml from a byte stream
func (project *Project) from_ConfYamlBytes(logger log.Log, yamlBytes []byte) bool {
	// parse the config file contents as a ConfSource_projectyaml object
	source := new(conf_Yaml)
	if err := yaml.Unmarshal(yamlBytes, source); err != nil {
		logger.Warning("YAML parsing error : " + err.Error())
		return false
	}
	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "YAML source:", *source)

	return source.configureProject(logger, project)
}

// A project configuration from Yaml
type conf_Yaml struct {
	Project string `yaml:"Project,omitempty"`
	Author  string `yaml:"Author,omitempty"`

	Paths map[string]string `yaml:"Paths,omitempty"`

	Tokens map[string]string `yaml:"Tokens,omitempty"`

	Settings map[string]string `yaml:"Settings,omitempty"`
}

// Make a Yaml Conf apply configuration to a project object
func (conf *conf_Yaml) configureProject(logger log.Log, project *Project) bool {
	// set a project name
	if conf.Project != "" {
		project.Name = conf.Project
	}
	// set a author name
	if conf.Author != "" {
		project.Author = conf.Author
	}

	// set any paths
	for key, keyPath := range conf.Paths {
		project.SetPath(key, keyPath, true)
	}

	// set any tokens
	for key, value := range conf.Tokens {
		project.SetToken(key, value)
	}

	/**
	 * Yaml Settings set Project Flags
	 */
	for key, value := range conf.Settings {
		switch key {
		case "UsePathsAsTokens":
			project.UsePathsAsTokens = conf.SettingStringToFlag(value)
		case "UseEnvVariablesAsTokens":
			project.UseEnvVariablesAsTokens = conf.SettingStringToFlag(value)
		}
	}

	logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Configured project from YAML conf")
	return true
}

func (conf *conf_Yaml) SettingStringToFlag(value string) bool {
	switch strings.ToLower(value) {
	case "y":
		fallthrough
	case "yes":
		fallthrough
	case "true":
		fallthrough
	case "1":
		return true

	default:
		return false
	}
}
