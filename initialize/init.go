package initialize

import (
	"github.com/james-nesbitt/coach-tools/conf"
	"github.com/james-nesbitt/coach-tools/log"
)

type InitTasks struct {
	conf *conf.Project

	root string

	tasks []InitTask
}

func (tasks *InitTasks) Init(logger log.Log, project *conf.Project, root string) bool {
	tasks.conf = project
	tasks.root = root
	tasks.tasks = []InitTask{}
	return true
}
func (tasks *InitTasks) RunTasks(logger log.Log) {
	for _, task := range tasks.tasks {
		logger.Debug(log.VERBOSITY_DEBUG_LOTS, "INIT TASK:", task)
		task.RunTask(logger)
	}
}

func (tasks *InitTasks) AddTask(task InitTask) {
	tasks.tasks = append(tasks.tasks, task)
}

func (tasks *InitTasks) AddFile(path string, contents string) {
	tasks.AddTask(InitTask(&InitTaskFile{
		root:     tasks.root,
		path:     path,
		contents: contents,
	}))
}
func (tasks *InitTasks) AddRemoteFile(path string, url string) {
	tasks.AddTask(InitTask(&InitTaskRemoteFile{
		root: tasks.root,
		path: path,
		url:  url,
	}))
}
func (tasks *InitTasks) AddFileCopy(path string, source string) {
	tasks.AddTask(InitTask(&InitTaskFileCopy{
		root:   tasks.root,
		path:   path,
		source: source,
	}))
}
func (tasks *InitTasks) AddGitClone(path string, url string) {
	tasks.AddTask(InitTask(&InitTaskGitClone{
		root: tasks.root,
		path: path,
		url:  url,
	}))
}
func (tasks *InitTasks) AddMessage(message string) {
	tasks.AddTask(InitTask(&InitTaskMessage{
		message: message,
	}))
}
func (tasks *InitTasks) AddError(error string) {
	tasks.AddTask(InitTask(&InitTaskError{
		error: error,
	}))
}

type InitTask interface {
	RunTask(logger log.Log) bool
}

type InitTaskError struct {
	error string
}

func (task *InitTaskError) RunTask(logger log.Log) bool {
	logger.Error(task.error)
	return true
}

type InitTaskMessage struct {
	message string
}

func (task *InitTaskMessage) RunTask(logger log.Log) bool {
	logger.Message(task.message)
	return true
}
