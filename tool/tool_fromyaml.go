package tool

import (
	"io/ioutil"

	"encoding/json"
	"gopkg.in/yaml.v2"

	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

const (
	COACH_TOOL_YAMLFILE = "tools.yml"
)

// Look for project Tooligurations inside the project Toolpaths
func (tools *Tools) from_ToolYaml(logger log.Log, project *conf.Project) {
	for _, yamlToolFilePath := range project.Paths.GetConfSubPaths(COACH_TOOL_YAMLFILE) {
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Looking for YAML Tool file: "+yamlToolFilePath)
		tools.from_ToolYamlFilePath(logger, project, yamlToolFilePath)
	}
}

// Try to Tooligure a project by parsing yaml from a Tool file
func (tools *Tools) from_ToolYamlFilePath(logger log.Log, project *conf.Project, yamlFilePath string) bool {
	// read the Toolig file
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Could not read a YAML file: "+err.Error())
		return false
	}

	if !tools.from_ToolYamlBytes(logger.MakeChild(yamlFilePath), project, yamlFile) {
		logger.Warning("YAML marshalling of the YAML Tool file failed [" + yamlFilePath + "]: " + err.Error())
		return false
	}
	return true
}

// Try to Tooligure a project by parsing yaml from a byte stream
func (tools *Tools) from_ToolYamlBytes(logger log.Log, project *conf.Project, yamlBytes []byte) bool {
	if project != nil {
		// token replace
		tokens := &project.Tokens
		yamlBytes = []byte(tokens.TokenReplace(string(yamlBytes)))
	}

	var yaml_tools map[string]map[string]interface{}
	err := yaml.Unmarshal(yamlBytes, &yaml_tools)
	if err != nil {
		logger.Warning("Could not parse tool yaml:" + err.Error())
		return false
	}

	for name, tool_struct := range yaml_tools {
		switch tool_struct["Type"] {
		case "script":
			json_tool, _ := json.Marshal(tool_struct)
			var scriptTool Tool_Script
			err := json.Unmarshal(json_tool, &scriptTool)
			if err != nil {
				logger.Warning("Couldn't process tool [" + name + "] :" + err.Error())
				continue
			}

			tool := Tool(&scriptTool)
			tool.Init(logger.MakeChild("SCRIPT:"+name), project)

			tools.SetTool(name, tool)
		}

	}
	return true
}
