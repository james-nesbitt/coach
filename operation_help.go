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
	helpOperationName := "help"
	helpOperationFlags := []string{}
	if len(operation.flags)>0 {
		helpOperationName = operation.flags[0]
	}
	if len(operation.flags)>1 {
		helpOperationFlags = operation.flags[1:]
	}

	helpOperation := GetOperation(helpOperationName, operation.nodes , operation.targets, operation.conf, operation.log)
	helpOperation.Help(helpOperationFlags)

}
