package main

import (
	"path"
	"io/ioutil"
)

type Operation_Tool struct {
	log Log

	conf *Conf

	nodes Nodes
	targets []string

	flags []string

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
	if operation.tools==nil {
		operation.tools = Tools{}
	}

	// load any user tools
	operation.ToolsFromYaml("usercoach", true)
	// load any project specific tools
	operation.ToolsFromYaml("projectcoach", true)

	if operation.tool=="" {
		operation.log.Error("No tool specified")
	} else if tool := operation.tools[operation.tool]; tool==nil {
		operation.log.Error("Specified tool not found")
	} else {
		tool.Run(operation.flags)
	}
}


func (operation *Operation_Tool) ToolsFromYaml(toolPathKey string, overwrite bool) bool {
	operation.log.Debug(LOG_SEVERITY_DEBUG_LOTS,"Updating Conf from YAML")

	if toolPath, ok := operation.conf.Path(toolPathKey); ok {
		// get the path to where the config file should be
		toolPath = path.Join(toolPath, "tools.yml")

		operation.log.Debug(LOG_SEVERITY_DEBUG_WOAH,"coach tool file:"+toolPath)

		// read the config file
		yamlFile, err := ioutil.ReadFile(toolPath)
		if err!=nil {
			operation.log.Debug(LOG_SEVERITY_DEBUG_LOTS,"Could not read the YAML file ["+toolPath+"]: "+err.Error())
			return false
		}

		// replace tokens in the yamlFile
		yamlFile = []byte( operation.conf.TokenReplace(string(yamlFile)) )
		operation.log.Debug(LOG_SEVERITY_DEBUG_STAAAP,"YAML (tokenized):"+ string(yamlFile))

		operation.tools.GetToolsFromYaml(operation.conf, operation.log, yamlFile, overwrite)
		return true
	}
	return false
}
