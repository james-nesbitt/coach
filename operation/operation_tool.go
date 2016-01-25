package operation

import (
	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
	"github.com/james-nesbitt/coach/tool"
)

type ToolOperation struct {
	log  log.Log
	conf *conf.Project

	root string // Path to root

	flags []string

	toolPaths []string

	tool  string
	tools *tool.Tools
}

func (operation *ToolOperation) Id() string {
	return "tool"
}
func (operation *ToolOperation) Flags(flags []string) bool {
	operation.tool = ""

	// first flag is the tool name
	if len(flags) > 0 {
		operation.tool = flags[0]
		flags = flags[1:]
	}

	operation.flags = flags
	return true
}
func (operation *ToolOperation) Help(topics []string) {
	operation.log.Message(`Operation: Tool

Coach will attempt to run an external tool, as listed in either the 
project or user tool.yml file.  This conceptually allows a script to
be run, with added ENV variables from the conf system, without having
to worry about the path to the script.

`)
}

func (operation *ToolOperation) Run(logger log.Log) bool {
	logger.Info("running tool operation")

	if operation.tools == nil {
		operation.tools = &tool.Tools{}

		// load tools from tool paths
		operation.tools.Init(logger, operation.conf)
	}

	if operation.tool == "" {
		operation.log.Error("No tool specified")
	} else if tool, ok := operation.tools.Tool(operation.tool); !ok {
		operation.log.Error("Specified tool not found: " + operation.tool)
	} else {
		tool.Run(operation.flags)
	}

	return true
}
