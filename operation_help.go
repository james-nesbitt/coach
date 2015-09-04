package main

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

The first topic passed in is assumed to be a help operation.
`)
}

func (operation *Operation_Help) Run() {
	operation.log.Message("running help operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:HELP")
	helpTopicName := "help"
	helpTopicFlags := []string{}
	if len(operation.flags)>0 {
		helpTopicName = operation.flags[0]
	}
	if len(operation.flags)>1 {
		helpTopicFlags = operation.flags[1:]
	}

	switch helpTopicName {
		case "cli":
			operation.Topic_CLI(helpTopicFlags)
		case "targets":
			fallthrough
		case "cli:targets":
			operation.Topic_CLI_Targets(helpTopicFlags)

		case "settings":
			operation.Topic_Settings(helpTopicFlags)
		case "tokens":
			fallthrough
		case "settings:tokens":
			operation.Topic_Settings_Tokens(helpTopicFlags)
		case "secrets":
			fallthrough
		case "settings:secrets":
			operation.Topic_Settings_Secrets(helpTopicFlags)

		default: //assume this is an operation call
			helpTopic := GetOperation(helpTopicName, operation.nodes , operation.targets, operation.conf, operation.log)
			helpTopic.Help(helpTopicFlags)
	}
}

func (operation *Operation_Help) Topic_CLI(flags []string) {
	operation.log.Note(`Topic: CLI

{CLI HELP}

`)
}
func (operation *Operation_Help) Topic_CLI_Targets(flags []string) {
	operation.log.Note(`Topic: CLI:Targets

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
`)
}


func (operation *Operation_Help) Topic_Settings_Tokens(flags []string) {
	operation.log.Note(`Topic: Settings:Tokens

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


func (operation *Operation_Help) Topic_Settings(flags []string) {
	operation.log.Note(`Topic: Settings

{SETTINGS HELP}

`)
}
func (operation *Operation_Help) Topic_Settings_Secrets(flags []string) {
	operation.log.Note(`Topic: Settings:Secrets

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
