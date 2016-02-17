package main

import (
	"os"

	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"

	"github.com/james-nesbitt/coach/libs"
	"github.com/james-nesbitt/coach/operation"
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

	logger = log.MakeCliLog("coach-cli", os.Stdout, verbosity)
	logger.Debug(log.VERBOSITY_DEBUG, "Reporting initialization", logger.Verbosity())

	workingDir, _ := os.Getwd()
	logger.Debug(log.VERBOSITY_DEBUG, "Working Directory", workingDir)

	project = conf.MakeCoachProject(logger.MakeChild("conf"), workingDir)
	logger.Debug(log.VERBOSITY_DEBUG, "Project configuration", *project)

	logger.Debug(log.VERBOSITY_DEBUG, "Finished initialization", nil)
}

func main() {

	if !(operationName == "init" || project.IsValid(logger.MakeChild("Sanity Check"))) {
		logger.Error("Coach project configuration is not processable.  Execution halted. [" + operationName + "]")
		return
	}

	logger.Debug(log.VERBOSITY_DEBUG, "Starting CLI Processing", nil)

	logger.Debug(log.VERBOSITY_DEBUG, "Creating client factories", nil)

	// get a list of client factories that we can use for nodes
	clientFactories := libs.MakeClientFactories(logger.MakeChild("client-factories"), project)
	logger.Debug(log.VERBOSITY_DEBUG, "Factories", *clientFactories)

	// Build our list of nodes
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

	/**
	 * convert the target string list from the arguments into a target set
	 */
	targets := nodes.Targets(logger.MakeChild("targets"), mainTargets)
	logger.Debug(log.VERBOSITY_DEBUG, "Sorted Targets", mainTargets, targets)

	/**
	 * Maybe we can come up with a better default operation?
	 *
	 * Various scenarios:
	 *
	 * A. all nodes are of the same type, then pick a decent default operation per node type.
	 *
	 */
	if operationName==operation.DEFAULT_OPERATION {

		// Check target node types
		targetType := ""
		for _, targetId := range targets.TargetOrder() {
			target, _ := targets.Target(targetId)
			node, _ := target.Node()
			if targetType=="" {
				targetType = node.Type()
			} else if targetType==node.Type() {

			} else {
				targetType = ""
				break
			}
		}
		switch targetType {
		case "command":
			// Command containers default to run
			operationName = "run"
		default:
			// default to info
			operationName = "info"
		}
	}

	/**
	 * Create an operation set
	 */
	operations := operation.MakeOperation(logger.MakeChild("operations"), project, operationName, operationFlags, targets)
	logger.Debug(log.VERBOSITY_DEBUG, "OPERATION:", operationName, operationFlags, operations)

	operations.Run(logger.MakeChild("operation"))

	logger.Debug(log.VERBOSITY_DEBUG, "Finished CLI Processing", nil)
}
