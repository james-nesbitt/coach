package main

import (
	"os"

	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
	"github.com/james-nesbitt/coach/operation"
)

var (
	globalFlags    map[string]string
	operationFlags []string

	logger  log.Log       // Logger interface for tracking messages
	project *conf.Project // project configuration
)

func init() {

	globalFlags, operationFlags = parseGlobalFlags(os.Args)

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

	operations := operation.MakeOperation(logger, project, "init", operationFlags, nil)
	operations.Run(logger)

}
