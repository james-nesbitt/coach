package operation

import (
	"strings"

	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/help"
	"github.com/james-nesbitt/coach/libs"
	"github.com/james-nesbitt/coach/log"
)

type HelpOperation struct {
	log  log.Log
	conf *conf.Project

	flags []string
}

func (operation *HelpOperation) Id() string {
	return "help"
}
func (operation *HelpOperation) Flags(flags []string) bool {
	operation.flags = flags
	return true
}
func (operation *HelpOperation) Help(topics []string) {
	operation.log.Message(`Operation: Help

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
func (operation *HelpOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: info")

	helpTopicName := "help"
	helpTopicFlags := []string{}

	if len(operation.flags) > 0 {
		helpTopicName = operation.flags[0]
	}
	if len(operation.flags) > 1 {
		helpTopicFlags = operation.flags[1:]
	}

	Helper := operation.getHelpObject()

	if topic, ok := Helper.Topic(helpTopicName, helpTopicFlags); ok {

		operation.log.Message(topic)
		return true

	} else {

		for _, operationName := range ListOperations() {
			if strings.HasPrefix(helpTopicName, operationName) {
				if helpOperations := MakeOperation(logger, operation.conf, operationName, operation.flags, &libs.Targets{}); len(helpOperations.operationsList) > 0 {
					for _, helpOperation := range helpOperations.operationsList {
						helpOperation.Help(append([]string{helpTopicName}, helpTopicFlags...))
					}
					return true
				}
			}
		}

	}

	operation.log.Warning("Unknown help topic")
	return false
}

// get a Help object, loaded with the default Help sources
func (operation *HelpOperation) getHelpObject() help.Help {
	// empty help object
	Helper := help.Help{}
	Helper.Init(operation.log, operation.conf)

	return Helper
}
