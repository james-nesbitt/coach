package libs

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

const (
	COACH_NODES_YAMLFILE       = "nodes.yml" // nodes are kept in the nodes.yml file
	NODES_YAML_DEFAULTNODETYPE = "service"   // by default we assume that a node is a service node
)

// Look for nodes configurations inside the project confpaths
func (nodes *Nodes) from_NodesYaml(logger log.Log, project *conf.Project, clientFactories *ClientFactories, overwrite bool) {
	for _, yamlNodesFilePath := range project.Paths.GetConfSubPaths(COACH_NODES_YAMLFILE) {
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Looking for YAML nodes file: "+yamlNodesFilePath)
		nodes.from_NodesYamlFilePath(logger, project, clientFactories, yamlNodesFilePath, overwrite)
	}
}

// Try to configure a project by parsing yaml from a conf file
func (nodes *Nodes) from_NodesYamlFilePath(logger log.Log, project *conf.Project, clientFactories *ClientFactories, yamlFilePath string, overwrite bool) bool {
	// read the config file
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Could not read a YAML file: "+err.Error())
		return false
	}

	if !nodes.from_NodesYamlBytes(logger.MakeChild(yamlFilePath), project, clientFactories, yamlFile, overwrite) {
		logger.Warning("YAML marshalling of the YAML nodes file failed [" + yamlFilePath + "]: " + err.Error())
		return false
	}
	return true
}

// Try to configure factories by parsing yaml from a byte stream
func (nodes *Nodes) from_NodesYamlBytes(logger log.Log, project *conf.Project, clientFactories *ClientFactories, yamlBytes []byte, overwrite bool) bool {
	if project != nil {
		// token replace
		tokens := &project.Tokens
		yamlBytes = []byte(tokens.TokenReplace(string(yamlBytes)))
	}

	var nodes_yaml map[string]node_yaml_v2
	err := yaml.Unmarshal(yamlBytes, &nodes_yaml)
	if err != nil {
		logger.Warning("YAML parsing error : " + err.Error())
		return false
	}
	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "YAML source:", nodes_yaml)

NodesListLoop:
	for name, node_yaml := range nodes_yaml {

		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Yaml Node:", name, node_yaml)

		nodeLogger := logger.MakeChild(name)

		if _, exists := nodes.Node(name); exists && !overwrite {
			nodeLogger.Warning("YAML node key already exists")
			continue NodesListLoop
		}
		if node_yaml.Disabled {
			nodeLogger.Warning("YAML node key is marked as disabled")
			continue NodesListLoop
		}

		var node Node

		// Start off assuming a default type
		nodeType := NODES_YAML_DEFAULTNODETYPE
		// if the node conf has a type, then use it
		if explicitType, ok := node_yaml.Type(); ok {
			nodeType = explicitType
		}

		switch strings.ToLower(strings.TrimSpace(nodeType)) {
		case "command":
			node = Node(&CommandNode{})
		case "build":
			node = Node(&BuildNode{})
		case "volume":
			node = Node(&VolumeNode{})
		case "service":
			node = Node(&ServiceNode{})
		default:
			nodeLogger.Warning("YAML node is an unknown type: " + nodeType)
			continue NodesListLoop
		}

		if node != nil {
			var client Client
			var instancesSettings InstancesSettings

			var ok bool

			if client, ok = node_yaml.GetClient(nodeLogger, clientFactories); !ok {
				nodeLogger.Error("Invalid Client configuration in node ")
			}
			if instancesSettings, ok = node_yaml.GetInstancesSettings(nodeLogger); !ok {
				nodeLogger.Error("Invalid Instances configuration in node")
			}

			nodeLogger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Building node from components:", client, instancesSettings)

			// run the node initializaer with the collected obkjects
			if node.Init(nodeLogger, name, project, client, instancesSettings) {
				// add any manual dependencies specified
				for _, dependency := range node_yaml.Requires {
					node.AddDependency(dependency)
				}

				nodeLogger.Debug(log.VERBOSITY_DEBUG_LOTS, "Adding node to nodes list:", name, node)
				nodes.SetNode(name, node, true)
			}
		}

	}

	return true
}

// 1. Original V1 coach nodes yaml format
type node_yaml_v1 struct {
	Disabled bool   `yaml:"Disabled,omitempty"`
	Type     string `yaml:"Type,omitempty"`

	Build string `yaml:"Build,omitempty"`
	Tag   string `yaml:"RepoTag,omitempty"`

	Instances string `yaml:"Instances,omitempty"`

	Docker FSouza_ClientSettings `yaml:"Docker,omitempty"`

	Requires []string `yaml:"Requires,omitempty"`
}

// 2. V2 Coach yaml format, with fixed fields
type node_yaml_v2 struct {
	Disabled bool   `yaml:"Disabled,omitempty"`
	NodeType string `yaml:"Type,omitempty"`

	ScaledInstances ScaledInstancesSettings `yaml:"Scale,omitempty"`
	FixedInstances  FixedInstancesSettings  `yaml:"Instances,omitempty"`
	SingleInstances bool                    `yaml:"Single,omitempty"`
	TempInstances   bool                    `yaml:"Disposable,omitempty"`

	Docker FSouza_ClientSettings `yaml:"Docker,omitempty"`

	Requires []string `yaml:"Requires,omitempty"`
}

func (node *node_yaml_v2) Type() (string, bool) {
	return node.NodeType, node.NodeType != ""
}
func (node *node_yaml_v2) GetClient(logger log.Log, clientFactories *ClientFactories) (Client, bool) {

	// if a docker client was configured then try to take it.
	// if !(node.Docker.Config.Image=="" && node.Docker.BuildPath=="") {
	if factory, ok := clientFactories.MatchClientFactory(FactoryMatchRequirements{Type: "docker"}); ok {
		if client, ok := factory.MakeClient(logger, ClientSettings(&node.Docker)); ok {
			return client, true
		}
	} else {
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Failed to match client factory:", factory)
	}
	// }

	logger.Warning("Invalid YAML node settings: improper client configuration")
	return nil, false
}
func (node *node_yaml_v2) GetInstancesSettings(logger log.Log) (InstancesSettings, bool) {
	var instancesSettings InstancesSettings

	if node.ScaledInstances.Maximum > 0 {
		instancesSettings = InstancesSettings(&node.ScaledInstances)
	} else if len([]string(node.FixedInstances.Names)) > 0 {
		instancesSettings = InstancesSettings(&node.FixedInstances)
	} else if bool(node.TempInstances) {
		instancesSettings = InstancesSettings(&TemporaryInstancesSettings{"run"})
	} else if bool(node.SingleInstances) {
		instancesSettings = InstancesSettings(&SingleInstancesSettings{Name: "single"})
	} else {
		instancesSettings = InstancesSettings(&NullInstancesSettings{})
	}

	return instancesSettings, true
}

// 3. Dynamic map based format for yaml nodes
type node_yaml_interface map[string]interface{}
