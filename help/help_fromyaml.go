package help

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

const (
	COACH_HELP_YAMLFILE = "help.yml"
)

// Ummarshaller interface : pass incoming yaml into the topics
func (help *Help) UnmarshalYAML(unmarshal func(interface{}) error) error {
	help.log.Debug(log.VERBOSITY_MESSAGE, "UNMARSHALL:", help.topics)
	return unmarshal(&help.topics)
}

// Look for help from the default core help, which is kept as yaml
func (help *Help) from_CoreHelpYaml(logger log.Log) {
	help.from_HelpYamlBytes(logger, nil, help.getCoreHelpYaml())
}

// Look for help inside the project confpaths
func (help *Help) from_HelpYaml(logger log.Log, project *conf.Project) {
	for _, yamlHelpFilePath := range project.Paths.GetConfSubPaths(COACH_HELP_YAMLFILE) {
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Looking for YAML Help file: "+yamlHelpFilePath)
		help.from_HelpYamlFilePath(logger, project, yamlHelpFilePath)
	}
}

// Try to configure help by parsing yaml from a Help file
func (help *Help) from_HelpYamlFilePath(logger log.Log, project *conf.Project, yamlFilePath string) bool {
	// read the Helpig file
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Could not read a YAML file: "+err.Error())
		return false
	}

	if !help.from_HelpYamlBytes(logger.MakeChild(yamlFilePath), project, yamlFile) {
		logger.Warning("YAML marshalling of the YAML Help file failed [" + yamlFilePath + "]: " + err.Error())
		return false
	}
	return true
}

// Try to configure help by parsing yaml from a byte stream
func (help *Help) from_HelpYamlBytes(logger log.Log, project *conf.Project, yamlBytes []byte) bool {
	if project != nil {
		// token replace
		tokens := &project.Tokens
		yamlBytes = []byte(tokens.TokenReplace(string(yamlBytes)))
		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Tokenized Bytes", string(yamlBytes))
	}

	if err := yaml.Unmarshal(yamlBytes, help); err != nil {
		logger.Warning("YAML parsing error : " + err.Error())
		// logger.Debug(log.VERBOSITY_DEBUG, "YAML parsing error : " + err.Error(), string(yamlBytes))
		return false
	}
	return true
}
