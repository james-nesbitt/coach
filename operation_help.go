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

TOPICS:
  cli
    cli:targets (targets)

	settings
	  settings:targets (targets)
	  settings:secrets (secrets)

OPERATIONS:

  init

	build
	pull
	destroy

	create
	remove

	start
	pause
	unpause
	remove

	scale

	tool

The first topic passed in is assumed to be a help operation.
`)
}

func (operation *Operation_Help) Run() {
	operation.log.Message("running help operation")
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

for name, _ := range help.Topics {
	operation.log.Debug(LOG_SEVERITY_MESSAGE, "TOPIC:"+name)
}

	if topic, ok := help.getTopic( helpTopicName, helpTopicFlags ); ok {

		operation.log.Note( topic )

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
help: |

  Coach will attempt to output help messages.  The message will match either a topic, or an operation.

  TOPICS:
    cli
      cli:targets (targets)

    settings
      settings:targets (targets)
      settings:secrets (secrets)

  OPERATIONS:

    init

    build
    pull
    destroy

    create
    remove

    start
    pause
    unpause
    remove

  The first topic passed in is assumed to be a help operation.


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

  {SETTINGS HELP}

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

`)
} 