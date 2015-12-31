package main

import (
	"os"

	"github.com/james-nesbitt/coach-tools/conf"
	"github.com/james-nesbitt/coach-tools/log"

	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/operation"
)

var (
	operationName  string
	mainTargets    []string
	globalFlags    map[string]string
	operationFlags []string

	logger  log.Log       // Logger interface for tracking messages
	project *conf.Project // project configuration
)

func init() {

	operationName, mainTargets, globalFlags, operationFlags = parseGlobalFlags(os.Args)

	// verbosity
	var verbosity int = log.VERBOSITY_MESSAGE
	if globalFlags["verbosity"] != "" {
		switch globalFlags["verbosity"] {
		case "message":
			verbosity = log.VERBOSITY_MESSAGE
		case "info":
			verbosity = log.VERBOSITY_INFO
		case "warning":
			verbosity = log.VERBOSITY_WARNING
		case "verbose":
			verbosity = log.VERBOSITY_DEBUG_LOTS
		case "debug":
			verbosity = log.VERBOSITY_DEBUG_WOAH
		case "staaap":
			verbosity = log.VERBOSITY_DEBUG_STAAAP
		}
	}

	logger = log.MakeCoachLog("coach", os.Stdout, verbosity)
	logger.Debug(log.VERBOSITY_DEBUG, "Reporting initialization", logger.Verbosity())

	workingDir, _ := os.Getwd()
	logger.Debug(log.VERBOSITY_DEBUG, "Working Directory", workingDir)

	project = conf.MakeCoachProject(logger.MakeChild("conf"), workingDir)
	logger.Debug(log.VERBOSITY_DEBUG, "Project configuration", *project)

	logger.Debug(log.VERBOSITY_DEBUG, "Finished initialization", nil)
}

func main() {

	if !project.IsValid(logger.MakeChild("Sanity Check")) {
		logger.Error("Coach project configuration is not processable.  Execution halted.")
		return
	}

	logger.Debug(log.VERBOSITY_DEBUG, "Starting CLI Processing", nil)

	logger.Debug(log.VERBOSITY_DEBUG, "Creating client factories", nil)

	clientFactories := libs.MakeClientFactories(logger.MakeChild("client-factories"), project)
	logger.Debug(log.VERBOSITY_DEBUG, "Factories", *clientFactories)

	nodes := libs.MakeNodes(logger.MakeChild("nodes"), project, clientFactories)
	logger.Debug(log.VERBOSITY_DEBUG, "Nodes", *nodes)

	/**
	 * prepare the whole node set
	 *
	 * this means that the nodes building is completed, and now any 
	 * dependencies and settings should already be included. The prepare
	 * should now connect pieces together down the chain.
	 */
	nodes.Prepare(logger.MakeChild("nodes"))

for name, node := range nodes.NodesMap {
	nodeLogger := logger.MakeChild(name)

	nodeLogger.Debug(log.VERBOSITY_DEBUG_STAAAP, "node", node)
	nodeLogger.Debug(log.VERBOSITY_DEBUG_STAAAP, "instances", node.Instances())
	nodeLogger.Debug(log.VERBOSITY_DEBUG_STAAAP, "client", node.NodeClient())
}


	logger.Debug(log.VERBOSITY_DEBUG, "OPERATIONNAME:", operationName)
	logger.Debug(log.VERBOSITY_DEBUG, "TARGETS:", mainTargets)
	logger.Debug(log.VERBOSITY_DEBUG, "FLAGS:GLOBAL:", globalFlags)
	logger.Debug(log.VERBOSITY_DEBUG, "FLAGS:OPERATION:", operationFlags)

	logger.Debug(log.VERBOSITY_DEBUG, "LOG:", logger)

	logger.Debug(log.VERBOSITY_DEBUG, "Finished CLI Processing", true)
}

/**
 * Parse command flags to configure the operation
 *
 * 1: GLOBAL FLAGS : only those which we recognize below
 * 2: OPERATION [optional] : the first non-global flag, if we recognize it
 * 3. OPERATION ARGUMENTS : anything left
 *
 */
func parseGlobalFlags(flags []string) (operationName string, targetIdentifiers []string, globalFlags map[string]string, operationFlags []string) {
	operationName = operation.DEFAULT_OPERATION // default operation, to be interpreted later, if not set in this function

	globalFlags = map[string]string{} // start of with no flags
	targetIdentifiers = []string{}    //  ||

	global := true // start of assuming everything is a global arg
	for index := 1; index < len(flags); index++ {
		arg := flags[index]

		switch arg {
		case "-v":
			fallthrough
		case "--info":
			globalFlags["verbosity"] = "info"
		case "-vv":
			fallthrough
		case "--verbose":
			globalFlags["verbosity"] = "verbose"
		case "-vvv":
			fallthrough
		case "--debug":
			globalFlags["verbosity"] = "debug"
		case "-vvvv":
			fallthrough
		case "--staaap":
			globalFlags["verbosity"] = "staaap"
			fallthrough

		case "--all": // this is default anyway
			targetIdentifiers = append(targetIdentifiers, "$all")

		default:

			/**
			* The first flags that we don't recognize as global, fall into three cases:
			*  @{flag} : indicates a node target, can be repeated
			*  %{flag} : indicates a node type target, can be repeated
			*  -{flag} : indicates the end of global flag targeting, and starts the collection of operationFlags
			*  {flag} : (first only) indicates which operation (default is info)
			 */

			switch arg[0:1] {
			case "@": // target
				fallthrough
			case "%": // type
				targetIdentifiers = append(targetIdentifiers, arg)

			// this means that local flags have started being processed, as all global flags are particular
			case "-": // local flag
				global = false

			default: // operation

				// if we recognize the next argument as an operation, then set it,
				// otherwise we assume a default operation, and that op args have started
				// @TODO there has got to be a better way of doing this.
				if operation.IsValidOperationName(arg) {
					operationName = arg
					index++
				}

				global = false

			}

		}

		// all remaining flags are local
		if !global {
			operationFlags = flags[index:]
			break
		}
	}

	// return is handles via named arguments
	return
}
