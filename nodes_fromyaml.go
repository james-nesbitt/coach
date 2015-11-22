package main

import (
	"path"
	"strings"
	"strconv"

	"io/ioutil"

	"gopkg.in/yaml.v2"
	docker "github.com/fsouza/go-dockerclient"
)

const nodesource_nodesyaml_filename string = "nodes.yml"

func (nodes *Nodes) from_Yaml(log Log, conf *Conf) {
	log.Debug(LOG_SEVERITY_DEBUG_LOTS,"Updating Conf from YAML")

	if nodesPath, ok := conf.Path("projectcoach"); ok {
		// get the path to where the config file should be
		nodesPath = path.Join(nodesPath, nodesource_nodesyaml_filename)

		log.Debug(LOG_SEVERITY_DEBUG_WOAH,"Project coach nodes file:"+nodesPath)

		// read the config file
		yamlFile, err := ioutil.ReadFile(nodesPath)
		if err!=nil {
			log.Warning("Could not read the YAML file: "+err.Error())
			return
		}

		// replace tokens in the yamlFile
		yamlFile = []byte( conf.TokenReplace(string(yamlFile)) )
		log.Debug(LOG_SEVERITY_DEBUG_STAAAP,"YAML (tokenized):"+ string(yamlFile))

		// parse the config file contents as a ConfSource_projectyaml object
		source := new(Nodes_Yaml)
		if err := yaml.Unmarshal(yamlFile, source); err!=nil {
			log.Error("YAML marshalling of the YAML conf file failed: "+err.Error())
			return
		}
		log.DebugObject(LOG_SEVERITY_DEBUG_STAAAP,"YAML source object:", *source)

		for name, item := range *source {
			log.DebugObject(LOG_SEVERITY_DEBUG_STAAAP,"YAML ITEM:", name, "::", item)
			if node,err := item.toNode(name, log.ChildLog("NODE:"+name), conf); err==nil {
				nodes.AddNode(name, node)

				log.DebugObject(LOG_SEVERITY_DEBUG_WOAH,"YAML Node:", name, node)
			} else {
				log.Error(err.Error())
			}
		}

	} else {
		log.Debug(LOG_SEVERITY_DEBUG_LOTS,"YAML file not found, no project coach folder:"+nodesPath)
	}
}

type Nodes_Yaml map[string]Node_Yaml

type Node_Yaml struct {
	Disabled bool										`yaml:"Disabled,omitempty"`

	Type string											`yaml:"Type,omitempty"`

	Build string										`yaml:"Build,omitempty"`
	Tag string											`yaml:"RepoTag,omitempty"`

	Instances string								`yaml:"Instances,omitempty"`

	Config docker.Config						`yaml:"Config,omitempty"`
	Host docker.HostConfig					`yaml:"Host,omitempty"`

	Requires []string								`yaml:"Requires,omitempty"`
}

/**
 * Convert this YAML Node to a proper Node
 *
 * First we create a pretty simple node to pass in the Docker entities,
 * and then we run the Node Init() method to fill out the missing fields,
 * and then we run the instance builder appropriate for this .Type, which
 * should create appropriate instance objects (although they wont be processed
 * yet)
 */
func (item *Node_Yaml) toNode(name string, log Log, conf *Conf) (*Node, error) {

	if item.Type=="" {
		// there was not Type in the yaml for this node, so assume it is a service
		item.Type = "service"
	}

	node := Node{
		conf: conf,
		log: log,

		Name: name,

		BuildPath: item.Build,

		Config: item.Config,
		HostConfig: item.Host,
	}

	node.Init(item.Type)

	// interpret the Instances value
	if item.Instances=="" {
		if item.Type=="command" {
			item.Instances="temporary" // command nodes default to temporary
		} else {
			item.Instances="single" // all other nodes default to single
		}
	}

	// convert the instances string to proper node instances, using the first value as an optional instance type
	split := strings.Split(item.Instances, " ")
	switch split[0] {
		case "scaled":
			fallthrough
		case "scale":
			var min, max int64
			if len(split)<2 {
				min = 3
			} else {
				min, _ = strconv.ParseInt(split[1], 0, 0)
			}
			if len(split)<3 {
				max = 6
			} else {
				max, _ = strconv.ParseInt(split[2], 0, 0)
			}
			node.ConfigureInstances_Scaled( int(min), int(max))

		case "temp":
			fallthrough
		case "temporary":
			node.ConfigureInstances_Temporary()

		case "single":
			node.ConfigureInstances_Single()

		case "fixed":
			if len(split)<2 {
				split = []string{"unnamed"} // this was not properly configured, so just give a single instance name
			} else {
				split = split[1:] // remove the "fixed" and assume that the rest are instance names
			}
			fallthrough
		default: // default is a fixed type, space delimited list of instances
			node.ConfigureInstances_Fixed( split )
	}

	log.DebugObject(LOG_SEVERITY_DEBUG_STAAAP,"YAML item to node:", node)

	return &node, nil
}
