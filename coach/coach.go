package main

import (

	"github.com/james-nesbitt/coach-tools/client"
	"github.com/james-nesbitt/coach-tools/project"
	"github.com/james-nesbitt/coach-tools/log"
	"github.com/james-nesbitt/coach-tools/node"

	"github.com/james-nesbitt/coach-tools/operation"
)

const (
	USER_COACH_SUBFOLDER = ".coach" // most apps these days use ~/.config/
  PROJECT_COACH_SUBFOLDER = ".coach" // this subfolder marks the root of a coach project
)

var (
	operationName  string
	mainTargets    []string
	globalFlags    map[string]string
	operationFlags []string

	mainLog    log.Log // Logger interface for tracking messages
	mainProject   *project.Project
	mainClient *client.Client
	mainNodes  *node.Nodes
)

func init() {

	operationName, mainTargets, globalFlags, operationFlags = parseGlobalFlags(os.Args)

	// verbosity
	var verbosity int = log.VERBOSITY_MESSAGE
	if globalFlags["verbosity"]!="" {
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

	mainLog = log.GetLog("coach", verbosity)

  mainLog.Debug(log.VERBOSITY_DEBUG, "Reporting initialization", nil)

  mainLog.Debug(log.VERBOSITY_DEBUG, "Main Logger generated", mainLog.Verbosity())

	mainLog.Debug(log.VERBOSITY_DEBUG, "Generating main Project confiugration")
	mainProject = conf.GetProject(mainLog.MakeChild("Project"))
  mainLog.Debug(log.VERBOSITY_DEBUG, "Main Project configuration", *mainProject)

	mainLog.Debug(log.VERBOSITY_DEBUG, "Generating main Client")
	mainClient = client.GetClient()
  mainLog.Debug(log.VERBOSITY_DEBUG, "Main Client", *mainClient)

	mainLog.Debug(log.VERBOSITY_DEBUG, "Collecting Node List")
	mainNodes = node.GetNodes()
  mainLog.Debug(log.VERBOSITY_DEBUG, "Main Nodes list", *mainNodes)

  mainLog.Debug(log.VERBOSITY_DEBUG, "Finished initialization", nil)
}

func main() {

  if !mainProject.Valid() {
  	mainLog.Critical("Coach configuration is not processable.  Execution halted.")
  	return
	}

  mainLog.Debug(log.VERBOSITY_DEBUG, "Starting CLI Processing", nil)

	fmt.Println("OPERATIONNAME:", operationName)
	fmt.Println("TARGETS:", mainTargets)
	fmt.Println("FALGS:GLOBAL:", globalFlags)
	fmt.Println("FALGS:OPERATION:", operationFlags)




	fmt.Println("LOG:", mainLog)
	fmt.Println("PROJECT:", mainProject)
	fmt.Println("CLIENT:", mainClient)
	fmt.Println("NODES:", mainNodes)

  mainLog.Debug(log.VERBOSITY_DEBUG, "Finished CLI Processing", true)
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
