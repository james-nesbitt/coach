package tool

import (
	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

// DB of Tools
type Tools map[string]Tool

func (tools *Tools) Init(logger log.Log, project *conf.Project) bool {
	tools.from_ToolYaml(logger, project)
	return true
}
func (tools Tools) Tool(name string) (tool Tool, ok bool) {
	tool, ok = tools[name]
	return
}
func (tools Tools) SetTool(name string, tool Tool) {
	tools[name] = tool
}

// Defining Tool interface
type Tool interface {
	Init(logger log.Log, project *conf.Project) bool
	Run(flags []string) bool
}
