package operation

import (
	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/libs"
	"github.com/james-nesbitt/coach/log"
)

const (
	DEFAULT_OPERATION = "<default operation>"
)

func MakeOperation(logger log.Log, project *conf.Project, name string, flags []string, targets *libs.Targets) *Operations {
	operations := Operations{}
	operations.Init(logger, OperationsSettings{}, targets)

	opLogger := logger.MakeChild(name)

	var operation Operation
	switch name {
	case "info":
		operation = Operation(&InfoOperation{log: opLogger, targets: targets})

	case "pull":
		operation = Operation(&PullOperation{log: opLogger, targets: targets})
	case "build":
		operation = Operation(&BuildOperation{log: opLogger, targets: targets})
	case "destroy":
		operation = Operation(&DestroyOperation{log: opLogger, targets: targets})

	case "create":
		operation = Operation(&CreateOperation{log: opLogger, targets: targets})
	case "remove":
		operation = Operation(&RemoveOperation{log: opLogger, targets: targets})
	case "start":
		operation = Operation(&StartOperation{log: opLogger, targets: targets})
	case "stop":
		operation = Operation(&StopOperation{log: opLogger, targets: targets})
	case "scale":
		operation = Operation(&ScaleOperation{log: opLogger, targets: targets})
	case "pause":
		operation = Operation(&PauseOperation{log: opLogger, targets: targets})
	case "unpause":
		operation = Operation(&UnpauseOperation{log: opLogger, targets: targets})

	case "commit":
		operation = Operation(&CommitOperation{log: opLogger, targets: targets})

	case "up":
		operation = Operation(&UpOperation{log: opLogger, targets: targets})
	case "clean":
		operation = Operation(&CleanOperation{log: opLogger, targets: targets})

	case "run":
		operation = Operation(&RunOperation{log: opLogger, targets: targets})

	case "help":
		operation = Operation(&HelpOperation{log: opLogger, conf: project})

	case "init":
		operation = Operation(&InitOperation{log: opLogger, conf: project})

	case "tool":
		operation = Operation(&ToolOperation{log: opLogger, conf: project})

	default:
		operation = Operation(&UnknownOperation{id: name})
	}
	operation.Flags(flags)

	operations.AddOperation(operation)

	return &operations
}

// Configure a set of operations
type OperationsSettings struct {
}

// A set of Operations
type Operations struct {
	log            log.Log
	settings       OperationsSettings
	targets        *libs.Targets
	operationsList []Operation
}

// Constructor for the operations object
func (operations *Operations) Init(logger log.Log, settings OperationsSettings, targets *libs.Targets) {
	operations.log = logger
	operations.settings = settings
	operations.targets = targets
	operations.operationsList = []Operation{}
}

// Add an operation to the list of operations to run.
func (operations *Operations) AddOperation(operation Operation) {
	operations.operationsList = append(operations.operationsList, operation)
}

// Run all of the prepared operations
func (operations *Operations) Run(logger log.Log) {
	if len(operations.operationsList) == 0 {
		logger.Error("No operation created")
	} else {
		for _, operation := range operations.operationsList {
			operation.Run(logger.MakeChild(operation.Id()))
		}
	}
}

// Operation that can act on A target list
type Operation interface {
	Id() string
	Flags(flags []string) bool
	Run(log.Log) bool
	Help(topics []string)
}

func ListOperations() []string {
	return []string{
		"help",
		"info",
		"init",
		"tool",
		"pull",
		"build",
		"clean",
		"destroy",
		"run",
		"up",
		"scale",
		"create",
		"remove",
		"start",
		"stop",
		"attach",
		"pause",
		"unpause",
		"commit",
	}
}

// Validate a string as an operation name
func IsValidOperationName(name string) bool {
	for _, operation := range ListOperations() {
		if operation == name {
			return true
		}
	}
	return false
}

/**
 * No operation found
 */
type UnknownOperation struct {
	id string
}

func (operation *UnknownOperation) Id() string {
	return "unknown"
}
func (operation *UnknownOperation) Flags(flags []string) bool {
	return true
}
func (operation *UnknownOperation) Help(flags []string) {
	return
}
func (operation *UnknownOperation) Run(logger log.Log) bool {
	logger.Error("Unknown operation: " + operation.id)
	return false
}
