package main

import (
	"os"
)

func main() {

	// Parse the command flag
	// operation, targets, globalFlags, operationFlags := parseGlobalFlags( os.Args )
	operationName, targets, globalFlags, operationFlags := parseGlobalFlags( os.Args )

	// Interpret some of the global vars

	// verbosity
	var verbosity int = LOG_SEVERITY_WARNING
	if globalFlags["verbosity"]!="" {
		switch globalFlags["verbosity"] {
			case "message":
				verbosity = LOG_SEVERITY_MESSAGE
			case "warning":
				verbosity = LOG_SEVERITY_WARNING
			case "verbose":
				verbosity = LOG_SEVERITY_DEBUG_LOTS
			case "debug":
				verbosity = LOG_SEVERITY_DEBUG_WOAH
			case "staaap":
				verbosity = LOG_SEVERITY_DEBUG_STAAAP
		}
	}

	/**
	 * LOG: Get a parent logging object to be used throughout execution
	 */

	log := GetLog(os.Stdout, verbosity)

	log.DebugObject(LOG_SEVERITY_DEBUG, "Global Flags:", globalFlags)
	log.DebugObject(LOG_SEVERITY_MESSAGE, "Operation Flags:", operationFlags)
	log.DebugObject(LOG_SEVERITY_DEBUG, "Initial Targets:", targets)

	/**
	 * CONF: Get a configuration object
	 */

	conf := GetConf( log.ChildLog("CONF") )

	log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "CONF", conf);

	/**
	 * DOCKER CLIENT: Get a docker client from the conf settings
	 */

	client, err := GetClient(conf, log.ChildLog("DOCKERCLIENT"))
	if err!=nil {
		log.Fatal("could not create a docker client.  Note that you may force docker settings for IP/CertPath in your conf.yml file")
		return
	}

	log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Docker Client", *client);

	/**
	 * NODES: Get the list of nodes that are a part of this project
	 */

	nodes := getNodes(log.ChildLog("NODES"), &conf, client, targets)

	log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "NODES MAP", nodes);

	/**
	 * Now we are fully configured and ready to go
	 */


	log.DebugObject(LOG_SEVERITY_DEBUG, "OPERATION: ["+operationName+"] => flags :", operationFlags)

	// get an operation object
	operation := GetOperation(operationName, nodes, targets, client, &conf, log.ChildLog("OPERATION"))

	log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Operation", operation);

	// pass in any flags
	operation.Flags(operationFlags)

	// Run the operation
	operation.Run()

}


/**
 * Parse command flags to configure the operation
 */
func parseGlobalFlags(flags []string) (operation string, targets []string, globalFlags map[string]string, operationFlags []string) {
	operation = "info"				// default operation

	globalFlags = map[string]string{}
	targets = []string{}

	global := true // start of assuming everything is a global arg
	for index:=1; index<len(flags); index++ {
		arg := flags[index]

		switch(arg) {
			case "-v":
				fallthrough
			case "--verbose":
				globalFlags["verbosity"] = "verbose"
			case "-vv":
				fallthrough
			case "--debug":
				globalFlags["verbosity"] = "debug"
			case "-vvv":
				fallthrough
			case "--staaap":
				globalFlags["verbosity"] = "staaap"
				fallthrough

			case "--all": // this is default anyway
				targets = append(targets, "$all")

			default:

				/**
				* The first flags that we don't recognize as global, fall into three cases:
				*  @{flag} : indicates a node target, can be repeated
				*  %{flag} : indicates a node type target, can be repeated
				*  -{flag} : indicates the end of global flag targeting, and starts the collection of operationFlags
				*  {flag} : (first only) indicates which operation (default is info)
				*/

				switch (arg[0:1]) {
					case "@": // target
						fallthrough
					case "%": // type
						targets = append(targets, arg)

					// this means that local flags have started being processed, as all global flags are particular
					case "-": // local flag
						global = false

					default: // operation
						operation = arg
						index++
						global = false
				}

		}

		// all remaining flags are local
		if (!global) {
			operationFlags = flags[index:]
			break
		}
	}

	// return is handles via named arguments
	return
}
