package main

import (
	"path"
)

type Operation_Tool struct {
	log Log

	conf *Conf

	nodes Nodes
	targets []string

	flags []string

	toolPaths []string

	tool string
	tools Tools
}
func (operation *Operation_Tool) Flags(flags []string) {
	operation.tool = ""

	// first flag is the tool name
	if len(flags)>0 {
		operation.tool = flags[0]
		flags = flags[1:]
	}

	operation.flags = flags

	operation.toolPaths = []string{"usercoach","projectcoach"}
}

func (operation *Operation_Tool) Help(topics []string) {
	operation.log.Note(`Operation: Tool

Coach will attempt to run an external tool, as listed in either the 
project or user tool.yml file.  This conceptually allows a script to
be run, with added ENV variables from the conf system, without having
to worry about the path to the script.


`)
}

func (operation *Operation_Tool) Run() {	
	operation.log.Info("running tool operation")	

	if operation.tools==nil {
		operation.tools = Tools{}
	}

	// load tools from tool paths
	for _, pathKey := range operation.toolPaths {
		operation.ToolsFromYaml(pathKey, true)
	}

	if operation.tool=="" {
		operation.log.Error("No tool specified")
	} else if tool := operation.tools[operation.tool]; tool==nil {
		operation.log.Error("Specified tool not found")
	} else {
		tool.Run(operation.flags)
	}
}

// populate the tools list from a yamlfile, in a conf path
func (operation *Operation_Tool) ToolsFromYaml(toolPathKey string, overwrite bool) bool {
	operation.log.Debug(LOG_SEVERITY_DEBUG_LOTS,"Updating Conf from YAML")

	if toolPath, ok := operation.conf.Path(toolPathKey); ok {
		// get the path to where the config file should be
		toolPath = path.Join(toolPath, "tools.yml")
		operation.tools.GetToolsFromYamlFile(operation.conf, operation.log, toolPath, overwrite)

		return true
	}
	return false
}
