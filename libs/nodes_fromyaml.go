package libs

import (
	"string"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/james-nesbitt/coach-tools/log"
	"github.com/james-nesbitt/coach-tools/conf"
)

const (
	COACH_NODES_YAMLFILE = "nodes.yml"
)

// Look for nodes configurations inside the project confpaths
func (nodes *Nodes) from_NodesYaml(logger log.Log, project *conf.Project, overwrite bool) {
	for _, yamlNodesFilePath := range project.Paths.GetConfSubPaths(COACH_NODES_YAMLFILE) {
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Looking for YAML nodes file: "+yamlNodesFilePath)
		nodes.from_NodesYamlFilePath(logger, yamlNodesFilePath, overwrite)
	}
}

// Try to configure a project by parsing yaml from a conf file
func (nodes *Nodes) from_NodesYamlFilePath(logger log.Log, yamlFilePath string, overwrite bool) bool {
	// read the config file
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Could not read a YAML file: "+err.Error())
		return false
	}

	if !nodes.from_NodesYamlBytes(logger.MakeChild(yamlFilePath), yamlFile, overwrite) {
		logger.Warning("YAML marshalling of the YAML nodes file failed [" + yamlFilePath + "]: " + err.Error())
		return false
	}
	return true
}

// Try to configure factories by parsing yaml from a byte stream
func (nodes *Nodes) from_NodesYamlBytes(logger log.Log, yamlBytes []byte, overwrite bool) bool {
	var nodes_yaml map[string]node_yaml_interface
	err := yaml.Unmarshal(yamlBytes, &nodes_yaml)
	if err != nil {
		logger.Warning("YAML parsing error : " + err.Error())
		return false
	}
	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "YAML source:", nodes_yaml)

	for name, node_yaml := range nodes_yaml {
		if _, exists := nodes.Node(name); overwrite || !exists {

			var node Node 
			nodeType := ""

			if yamlType, ok := node_yaml["Type"]; ok && yamlType.(string) {
				nodeType = yamlType
				delete(node_yaml, "Type")
			} else if yamlType, ok := node_yaml["type"]; ok && yamlType.(string) {}
				nodeType = yamlType
				delete(node_yaml, "type")
			}

			switch strings.ToLower(nodeType) {

			}


			for key, value := range node_yaml {

			}



		} else {
			logger.Warning("YAML node key already exists: "+name)
		}
	}

	return true
}



type node_yaml_v1 struct {
	Disabled bool   `yaml:"Disabled,omitempty"`
	Type     string `yaml:"Type,omitempty"`

	Build string `yaml:"Build,omitempty"`
	Tag   string `yaml:"RepoTag,omitempty"`

	Instances string `yaml:"Instances,omitempty"`

	Docker Client_DockerFSouzaSettings `yaml:"Docker,omitempty"`

	Requires []string `yaml:"Requires,omitempty"`
}

type node_yaml_v2 struct {
	Disabled bool   `yaml:"Disabled,omitempty"`
	Type     string `yaml:"Type,omitempty"`

	ScaledInstances instances_yaml_scaled `yaml:"Scaled,omitempty"`
	FixedInstances instances_yaml_fixed `yaml:"Instances,omitempty"`

	Docker Client_DockerFSouzaSettings `yaml:"Docker,omitempty"`

	Requires []string `yaml:"Requires,omitempty"`
}
type instances_yaml_scaled string
type instances_yaml_fixed []string

type node_yaml_interface map[string]interface{}
