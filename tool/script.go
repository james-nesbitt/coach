package tool

import (
	"os"
	"os/exec"
	"path"

	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"	
)

// Script type tool
type Tool_Script struct {
	conf *conf.Project
	log log.Log

	Script []string		`json:"Script,omitempty" yaml:"Script,omitempty"`
	Env []string		`json:"ENV,omitempty" yaml:"ENV,omitempty"`

	EnvIsolate bool     `json:"EnvIsolate,omitempty" yaml:"EnvIsolate,omitempty"`
	LeavePathAlone bool `json:"LeavePathAlone,omitempty" yaml:"LeavePathAlone,omitempty"`
}
func (tool *Tool_Script) Init(logger log.Log, project *conf.Project) bool {
	tool.log = logger
	tool.conf = project

	return true
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

	root := ""
	if path.IsAbs(cmd_first) {

	} else if cmd_first[0:1]=="~" { // one would think that path can handle such terminology
		root, _ = tool.conf.Paths.Path("user-home")
		cmd_first = cmd_first[1:]
	} else {
		root, _ = tool.conf.Paths.Path("project-root")
	}
	if cmd_first!= "" {
		cmd_first = path.Join(root, cmd_first)
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

	tool.log.Debug(log.VERBOSITY_INFO ,"SCRIPT RUN")
	err = cmd.Wait()
	tool.log.Message("FINISHED")
	return err==nil
}
