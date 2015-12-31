package libs

import (
	"strings"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/james-nesbitt/coach-tools/log"
	"github.com/james-nesbitt/coach-tools/conf"
)

const (
	COACH_NODES_YAMLFILE = "nodes.yml" // nodes are kept in the nodes.yml file
	NODES_YAML_DEFAULTNODETYPE = "service" // by default we assume that a node is a service node
)

// Look for nodes configurations inside the project confpaths
func (nodes *Nodes) from_NodesYaml(logger log.Log, project *conf.Project, clientFactories *ClientFactories, overwrite bool) {
	for _, yamlNodesFilePath := range project.Paths.GetConfSubPaths(COACH_NODES_YAMLFILE) {
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Looking for YAML nodes file: "+yamlNodesFilePath)
		nodes.from_NodesYamlFilePath(logger, clientFactories, yamlNodesFilePath, overwrite)
	}
}

// Try to configure a project by parsing yaml from a conf file
func (nodes *Nodes) from_NodesYamlFilePath(logger log.Log, clientFactories *ClientFactories, yamlFilePath string, overwrite bool) bool {
	// read the config file
	yamlFile, err := ioutil.ReadFile(yamlFilePath)
	if err != nil {
		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "Could not read a YAML file: "+err.Error())
		return false
	}

	if !nodes.from_NodesYamlBytes(logger.MakeChild(yamlFilePath), clientFactories, yamlFile, overwrite) {
		logger.Warning("YAML marshalling of the YAML nodes file failed [" + yamlFilePath + "]: " + err.Error())
		return false
	}
	return true
}

// Try to configure factories by parsing yaml from a byte stream
func (nodes *Nodes) from_NodesYamlBytes(logger log.Log, clientFactories *ClientFactories, yamlBytes []byte, overwrite bool) bool {
	var nodes_yaml map[string]node_yaml_v2
	err := yaml.Unmarshal(yamlBytes, &nodes_yaml)
	if err != nil {
		logger.Warning("YAML parsing error : " + err.Error())
		return false
	}
	logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "YAML source:", nodes_yaml)

	NodesListLoop:
		for name, node_yaml := range nodes_yaml {
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

			switch strings.ToLower(nodeType) {
			case "build":
				node = Node(&BuildNode{})
			case "volume":
				node = Node(&VolumeNode{})
			case "service":
				node = Node(&ServiceNode{})
			case "command":
				node = Node(&CommandNode{})
			default:
				nodeLogger.Warning("YAML node is an unknown type '"+nodeType)
				continue NodesListLoop
			}

			if node!=nil {
				var client Client
				var instances Instances

				var ok bool

				if client, ok = node_yaml.GetClient(nodeLogger, clientFactories); !ok {
					nodeLogger.Error("Invalid Client configuration in node ")
				}
				if instances, ok = node_yaml.GetInstances(nodeLogger, client); !ok {
					nodeLogger.Error("Invalid Instances configuration in node")
				}

				nodeLogger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Building node from components:", client, instances)

				// run the node initializaer with the collected obkjects
				if node.Init(nodeLogger, client, instances) {
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

	Docker Client_DockerFSouzaSettings `yaml:"Docker,omitempty"`

	Requires []string `yaml:"Requires,omitempty"`
}

// 2. V2 Coach yaml format, with fixed fields
type node_yaml_v2 struct {
	Disabled bool   `yaml:"Disabled,omitempty"`
	NodeType     string `yaml:"Type,omitempty"`

	ScaledInstances ScaledInstancesSettings `yaml:"Scale,omitempty"`
	FixedInstances FixedInstancesSettings `yaml:"Instances,omitempty"`
	TempInstances bool `yaml:"Disposable,omitempty"`

	Docker Client_DockerFSouzaSettings `yaml:"Docker,omitempty"`

	Requires map[string][]string `yaml:"Requires,omitempty"`
}

func (node *node_yaml_v2) Type() (string, bool) {
	return node.NodeType, node.NodeType!=""
}
func (node *node_yaml_v2) GetClient(logger log.Log, clientFactories *ClientFactories) (Client, bool) {

	// if a docker client was configured then try to take it.
	// if !(node.Docker.Config.Image=="" && node.Docker.BuildPath=="") {
		if factory, ok := clientFactories.MatchClientFactory( FactoryMatchRequirements{Type:"docker"} ); ok {
			if client, ok := factory.MakeClient(logger, ClientSettings(&node.Docker)); ok {
				return client, true
			}
		} else {
			logger.Debug(log.VERBOSITY_DEBUG_STAAAP,"Failed to match client factory:", factory)
		}
	// }

	logger.Warning("Invalid YAML node settings: improper client configuration")
	return nil, false
}
func (node *node_yaml_v2) GetInstances(logger log.Log, client Client) (Instances, bool) {

	var instancesSettings InstancesSettings
	var instances Instances

	if node.ScaledInstances.Maximum>0 {
		instancesSettings = InstancesSettings(&node.ScaledInstances)
		instances = Instances(&ScaledInstances{})
	} else if len([]string(node.FixedInstances.Names))>0 {
		instancesSettings = InstancesSettings(&node.FixedInstances)
		instances = Instances(&FixedInstances{})
	} else if bool(node.TempInstances) {
		instancesSettings = InstancesSettings(&TemporaryInstancesSettings{"run"})
		instances = Instances(&TemporaryInstances{})
	} else {
		instancesSettings = InstancesSettings(&SingleInstancesSettings{Name:"single"})
		instances = Instances(&SingleInstances{})
	}

	instances.Init(logger.MakeChild("instances"), client, instancesSettings)

	return instances, true
}

// 3. Dynamic map based format for yaml nodes
type node_yaml_interface map[string]interface{}
