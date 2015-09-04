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

Coach will attempt to output a help message.

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
		case "topics":
			operation.Topic_Targets(helpTopicFlags)

		default: //assume this is an operation call
			helpTopic := GetOperation(helpTopicName, operation.nodes , operation.targets, operation.conf, operation.log)
			helpTopic.Help(helpTopicFlags)
	}
}

func (operation *Operation_Help) Topic_Targets(flags []string) {
	operation.log.Note(`Topic: Targets

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
