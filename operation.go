package main

func GetOperation(name string, nodes Nodes, targets []string, conf *Conf, log Log) Operation {

	switch name {
		case "help":
			return Operation(&Operation_Help{log:log.ChildLog("HELP"), conf: conf, targets: targets})

		case "info":
			return Operation(&Operation_Info{log:log.ChildLog("INFO"), nodes:nodes, targets:targets})
		// case "status":

		case "init":
			return Operation(&Operation_Init{log:log.ChildLog("INIT"), conf: conf, targets: targets})

		case "pull":
			return Operation(&Operation_Pull{log:log.ChildLog("PULL"), nodes:nodes, targets:targets})
		case "build":
			return Operation(&Operation_Build{log:log.ChildLog("BUILD"), nodes:nodes, targets:targets})
		case "destroy":
			return Operation(&Operation_Destroy{log:log.ChildLog("DESTROY"), nodes:nodes, targets:targets})

		case "run":
			return Operation(&Operation_Run{log:log.ChildLog("RUN"), nodes:nodes, targets:targets})

		case "up":
			return Operation(&Operation_Up{log:log.ChildLog("UP"), nodes:nodes, targets:targets})

		case "create":
			return Operation(&Operation_Create{log:log.ChildLog("CREATE"), nodes:nodes, targets:targets})
		case "remove":
			return Operation(&Operation_Remove{log:log.ChildLog("REMOVE"), nodes:nodes, targets:targets})

		case "start":
			return Operation(&Operation_Start{log:log.ChildLog("START"), nodes:nodes, targets:targets})
		case "stop":
			return Operation(&Operation_Stop{log:log.ChildLog("STOP"), nodes:nodes, targets:targets, timeout: 3})

		case "attach":
			return Operation(&Operation_Attach{log:log.ChildLog("ATTACH"), nodes:nodes, targets:targets})

		case "pause":
			return Operation(&Operation_Pause{log:log.ChildLog("PAUSE"), nodes:nodes, targets:targets})
		case "unpause":
			return Operation(&Operation_Unpause{log:log.ChildLog("UNPAUSE"), nodes:nodes, targets:targets})

		case "commit":
			return Operation(&Operation_Commit{log:log.ChildLog("COMMIT"), nodes:nodes, targets:targets})

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
	operation.log.Note(`Topic: Missing

	No related help topic could be found.
`)
}
