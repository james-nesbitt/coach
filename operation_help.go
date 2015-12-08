package main

import (
	"path"
	"strings"

	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Operation_Help struct {
	log Log

	conf *Conf

	nodes Nodes
	targets []string

	flags []string
}
func (operation *Operation_Help) Flags(flags []string) {
	operation.flags = flags
}

func (operation *Operation_Help) Help(topics []string) {
	operation.log.Note(`Operation: HELP

Coach will attempt to output help messages.  The message will match either a topic, or an operation.

USAGE

	$/> coach help

	Default help page (this output)

	$/> coach help {topic}

	Help for a particular topic

	$/> coach help {operation}

	Help on a particular operation

TOPICS:

	cli: get help on how to use the cli

		cli:targets (targets) : get help about how targets work

	settings : get help about coach configuration

		settings;conf (conf) : get help on how to configure coach
		settings:nodes (nodes) : get help about how to define nodes
		settings:secrets (secrets) : get help about how to define secret tokens

OPERATIONS :

Target Independent: these operations don't pay attention to targets

  init: create a new coach project in the current path

	tool: run a project or user defined tool (see help tool)

Target Dependent: these operations will only act on passed targets	

  info: get information about project nodes

	pull: pull any node images
	build: build any node build images
	destroy: destroy any built node images

	create: create any needed node instance containers
	remove: remove any created node instance containers

	start: start node instances
	pause: pause all processes inside node instances
	unpause: pause all processes inside node instances
	remove: remove node instances (containers)

	scale: start (or stop) additional individual node instances to scale the app

	up: a shortcut operation for: build, pull, create, start
	clean: a shortcut operation for: stop, remove, destroy

The first topic passed in is assumed to be a help operation.
`)
}

func (operation *Operation_Help) Run() {
	operation.log.Info("running help operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)


	helpTopicName := "help"
	helpTopicFlags := []string{}
	if len(operation.flags)>0 {
		helpTopicName = operation.flags[0]
	}
	if len(operation.flags)>1 {
		helpTopicFlags = operation.flags[1:]
	}

	help := operation.getHelpObject()

	if topic, ok := help.getTopic( helpTopicName, helpTopicFlags ); ok {

		operation.log.Note( topic )
		return

	} else {

		operations := Operations{}
		for _, operationName := range operations.ListOperations() {
			if strings.HasPrefix(helpTopicName, operationName) {
				if helpOperation, ok := operations.GetOperation(operationName, operation.nodes, operation.targets, operation.conf, operation.log); ok {

					helpOperation.Help(append([]string{helpTopicName}, helpTopicFlags...))
					return

				}
			}
		}

	}

	operation.log.Warning("Unknown help topic")

}

// get a Help object, loaded with the default Help sources
func (operation *Operation_Help) getHelpObject() Help {
	// empty help object
  help := Help{conf: operation.conf, log: operation.log.ChildLog("Help")}

  // Add the core help
  help.GetHelpFromYaml( operation.getCoreHelpYaml() )

  // Look for help in a few config paths
  for _, helpPathKey := range []string{"projectcoach", "usercoach"} {

  	if helpPath, ok := operation.conf.Path(helpPathKey); ok {

  		helpPath = path.Join(helpPath, "help.yml")
			operation.log.Debug(LOG_SEVERITY_DEBUG_WOAH,"coach help file path:"+helpPath)
			help.HelpFromYamlFile( helpPath )

		}

	}

	return help
}


type Help struct {
	conf *Conf
	log Log

	Topics map[string]string   `yaml:"Topics,omitempty"`
}
// Ummarshaller interface : pass incoming yaml into the topics
func (help *Help) UnmarshalYAML(unmarshal func(interface{}) error) error {
help.log.DebugObject(LOG_SEVERITY_MESSAGE, "UNMARSHALL:",help.Topics)
	return unmarshal( &help.Topics )
}
func (help *Help) HelpFromYamlFile(helpPath string) bool {
	help.log.Debug(LOG_SEVERITY_DEBUG_LOTS,"Updating Help from YAML")

	// read the config file
	yamlFile, err := ioutil.ReadFile(helpPath)
	if err!=nil {
		help.log.Debug(LOG_SEVERITY_DEBUG_LOTS,"Could not read the YAML file ["+helpPath+"]: "+err.Error())
		return false
	}

	// replace tokens in the yamlFile
	yamlFile = []byte( help.conf.TokenReplace(string(yamlFile)) )
	help.log.Debug(LOG_SEVERITY_DEBUG_STAAAP,"YAML (tokenized):"+ string(yamlFile))

	help.GetHelpFromYaml(yamlFile)
	return true

}
func (help *Help) GetHelpFromYaml(source []byte) {
	merge := Help{}

	err := yaml.Unmarshal(source, &merge)
	if err!=nil {
		return
	}

	help.merge(merge)
}
func (help *Help) merge(merge Help) {
	if help.Topics==nil {
		help.Topics = map[string]string{}
	}
  for key, topic := range merge.Topics {
  	if _, ok := help.Topics[key]; !ok {
  		help.Topics[key] = topic
  	}
  }
}
func (help *Help) getTopic(topic string, flags []string) (string, bool) {
	if topic, ok := help.Topics[topic]; ok {
		return topic, true
	}

	return "", false
}

func (operation *Operation_Help) getCoreHelpYaml() []byte {
	return []byte(`

cli: |

	The CLI is the primary means of running coach.  The goal is typically to run a coach operation, on a number of coach nodes, or coach node instances.

	SEE ALSO:
	- cli:targets : $/> coach help cli:targets
	- operations : $/> coach help operations

"cli:targets": |

	Targets are a global setting used to determine which node and/or node instances an operation should used.  Targets are strings that define a type of node, a particular node, or a particular node instance.

	Coach accepts as a global flag, a list of targets in the following form:

	%{type} : all nodes of a certain type.  E.g.  %command
	%{type}:{instance} : the {instance} of any nodes of type {type}

	@{node} : all instances of a node named {node}
	@{node} : a particular {instance} instance from a node named {node}

	Here are some examples:

		$/> coach @db start
		Start all of the "db" node instances

		$/> coach @www.1 @www.2 remove
		remove the "1" and "2" instances from the "www" node

		$/> coach %service stop
		stop all nodes of type "service"

		$/> coach %volume:single commit
		commit the "single" instance of all nodes of type "volume"

settings: |
	Settings are primarily managed through a set of YAML files, that can be found in the project .coach folder.  In 
	the case of some settings files, copies can also exist in the users home folder at ~/.coach, in order to get
	user based configurations instead of project based.

	The most important settings concepts are:

	settings:conf : the general coach project configuration, found in the .coach/conf.yml file
	settings:nodes : the project nodes configuration, founc in the .coach/nodes.yml file 

	settings:tokens : read about how tokens can be defined in the conf.yml file, and used as tokens in the nodes.yml
	settings:secrets : additional sensitive tokens, that can be found in .coach/secrets/secrets.yml

	Note that other elements are often kept in the .coach file:

	- docker builds are often kept in Here
	- custom help topics ( settings:help:custom )

"settings:conf" : |
	The coach CLI has a number of file base configurations

	These configurations exist primarily in the .coach/conf.yml file.  More details can be found there
	or in the wiki.

	Typically the following settings are used:

	Project: define a string project name, which is used as the base name for built images and containers and is
		available as a token in the nodes.conf  %PROJECT

	Docker: you can override default docker settings, to for example explicitly use a remote docker client

	Tokens: a string map of token values, which can be user %token for string replacements in the nodes configuration file

"settings:nodes": |
	Nodes are any number of containers built around a single image, using similar configurations

	Nodes are defined in the .coach/nodes.yml file, and are used for the following purposes:

	- build: a node can define nothing more than a base docker build, used by other nodes
	- volume : a node can define a volume only container, which is created, but never started, and used to share volumes across an application
	- service : a node can define a service container, that is meant to be started and stopped, and possible scaled
	- command : a node can define a disposable command run container.

  More options are described in the wiki

"settings:tokens": |

	Tokens are settings key-value string pairs, that can be used for subsitution in later configurations.  Tokens defined in one state, are used for text replacement in all further configuration files, where possible.

	Tokens are typically used to:
		- allow settings to be reused across projects with minimal changes (which may allow templating in the future)
		- reduce duplicate settings across nodes, by centralizing values in the conf.yml
		- keep some sensitive values out of shared configuration (typically secrets)
		- allow some containers to use user specific values, instead of project specific values

	SOURCES:

	Tokens are typically kept in one of the following locations (list in order of loading):

	- .coach/conf.yml:Tokens: => in the conf.yml is a Tokens: map.  This map is typically used for values that are used across multiple nodes.

	- !/.coach/secrets/secrets.yml => in the user secrets.yml.  This map typically keeps user specific container ENV values such as passwords for user specific services that containers may use.
	- .coach/secrets/secrets.yml => in the project secrets.yml.  This map is typically used to keep project specific ID and token values used as ENV variables in containers, but that should not be kept in any source repository.

	For more information about secrets see $/> coach help secrets

"settings:secrets": |
	Secrets are just tokens, but they tend to be kept in locations that are easy to exclude from source versioning.

	For more information about tokens, see $/> coach help tokens

	User secrets allow users to have local values that are keyed for services that may be used in containers, that are not to be shared with other users, buy may be shared across projects.
	Project secrets keep secret tokens separate from conf tokens.  The project secrets could be tailored to each project user, or distributed separately from project source code.

	NOTE:
  	- tokens are not protected in any way during coach execution.

	SOURCES:

	Secrets are typically kept in one of the following locations (list in order of loading):

	- !/.coach/secrets/secrets.yml => in the user secrets.yml.  This map typically keeps user specific container ENV values such as passwords for user specific services that containers may use.
	- .coach/secrets/secrets.yml => in the project secrets.yml.  This map is typically used to keep project specific ID and token values used as ENV variables in containers, but that should not be kept in any source repository.


`)
} 