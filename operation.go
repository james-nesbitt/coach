package main

func GetOperation(name string, nodes Nodes, targets []string, conf *Conf, log Log) Operation {

	switch name {
		case "help":
			return Operation(&Operation_Help{log:log.ChildLog("HELP"), conf: conf, Targets: targets})

		case "info":
			return Operation(&Operation_Info{log:log.ChildLog("INFO"), Nodes:nodes, Targets:targets})
		// case "status":

		case "init":
			return Operation(&Operation_Init{log:log.ChildLog("INIT"), conf: conf, Targets: targets})

		case "pull":
			return Operation(&Operation_Pull{log:log.ChildLog("PULL"), Nodes:nodes, Targets:targets})
		case "build":
			return Operation(&Operation_Build{log:log.ChildLog("BUILD"), Nodes:nodes, Targets:targets})
		case "destroy":
			return Operation(&Operation_Destroy{log:log.ChildLog("DESTROY"), Nodes:nodes, Targets:targets})

		case "run":
			return Operation(&Operation_Run{log:log.ChildLog("RUN"), Nodes:nodes, Targets:targets})

		case "create":
			return Operation(&Operation_Create{log:log.ChildLog("CREATE"), Nodes:nodes, Targets:targets})
		case "remove":
			return Operation(&Operation_Remove{log:log.ChildLog("REMOVE"), Nodes:nodes, Targets:targets})

		case "start":
			return Operation(&Operation_Start{log:log.ChildLog("START"), Nodes:nodes, Targets:targets})
		case "stop":
			return Operation(&Operation_Stop{log:log.ChildLog("STOP"), Nodes:nodes, Targets:targets, timeout: 3})

		case "attach":
			return Operation(&Operation_Attach{log:log.ChildLog("ATTACH"), Nodes:nodes, Targets:targets})

		case "pause":
			return Operation(&Operation_Pause{log:log.ChildLog("PAUSE"), Nodes:nodes, Targets:targets})
		case "unpause":
			return Operation(&Operation_Unpause{log:log.ChildLog("UNPAUSE"), Nodes:nodes, Targets:targets})

		case "commit":
			return Operation(&Operation_Commit{log:log.ChildLog("COMMIT"), Nodes:nodes, Targets:targets})

		default:
			return Operation(&EmptyOperation{name: name, log: log.ChildLog("EMPTY")})
	}


	return nil
}


type Operation interface {
	Flags(flags []string)
	Run()

	Help(topics []string)
}


/**
 * @struc EmptyOperation A catchall operation if no other operation makes sense
 *
 * @TODO replace this with a help operation
 */

type EmptyOperation struct {
	name string
	log Log
}
func (operation *EmptyOperation) Flags(flags []string) {

}
func (operation *EmptyOperation) Run() {
	operation.log.Message("No matching operation found :"+operation.name)
}
func (operation *EmptyOperation) Help(topics []string) {
	operation.log.Note(`Operation: MissingOperation

	No related operation was found.
`)
}
