package main

import (
	"gopkg.in/yaml.v2"
	"encoding/json"

	"path"
	"os"
 	"os/exec"	
)

// DB of Tools
type Tools map[string]Tool

func (tools Tools) GetToolsFromYaml(conf *Conf, log Log, source []byte, overwrite bool) {
	var yaml_tools map[string]map[string]interface{}
	err := yaml.Unmarshal(source, &yaml_tools)
	if err!=nil {
		return
	}

	for name, tool_struct := range yaml_tools {

		switch (tool_struct["Type"]) {
			case "script":
				json_tool, _ := json.Marshal(tool_struct)
				var tool Tool_Script
				err := json.Unmarshal(json_tool, &tool)
				if err!=nil {
					continue
				}

				tool.conf = conf
				tool.log = log.ChildLog("SCRIPT")

				if exists:=tools[name]; exists==nil || overwrite {
					tools[name] = Tool(&tool)
				}
		}

	}
}

// Defining Tool interface
type Tool interface {
	Run(flags []string) bool
}

// Script type tool
type Tool_Script struct {
	conf *Conf
	log Log

	Script []string		`json:"Script,omitempty" yaml:"Script,omitempty"`
	Env []string		`json:"ENV,omitempty" yaml:"ENV,omitempty"`

	EnvIsolate bool     `json:"EnvIsolate,omitempty" yaml:"EnvIsolate,omitempty"`
	LeavePathAlone bool `json:"LeavePathAlone,omitempty" yaml:"LeavePathAlone,omitempty"`
}
func (tool *Tool_Script) Run(flags []string) bool {

	cmd_first := tool.Script[0]
	args := []string{}
	if len(tool.Script)>1 {
		args = tool.Script[1:]
	}
	if len(flags)>0 {
		args = append(args, flags...)
	}

	if path.IsAbs(cmd_first) {

	} else if cmd_first[0:1]=="~" { // one would think that path can handle such terminology
		cmd_first = path.Join(tool.conf.Paths["userhome"], cmd_first[1:])
	} else {
		cmd_first = path.Join(tool.conf.Paths["project"], cmd_first)
	}

	cmd := exec.Command( cmd_first, args... )
	cmd.Stdin = os.Stdin
	cmd.Stdout = tool.log
	cmd.Stderr = tool.log

	if len(tool.Env)>0 {
		if !tool.EnvIsolate {
			cmd.Env = append(cmd.Env, os.Environ()...)
		}

		cmd.Env = append(cmd.Env, tool.Env...) 
	}

	tool.log.Message("RUN: "+cmd_first)
	err := cmd.Start()

	if err!=nil {
		tool.log.Error("FAILED => "+err.Error())
		return false
	}

	err = cmd.Wait()
	tool.log.Message("FINISHED")
	return err==nil
}
