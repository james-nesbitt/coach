package main

import (
	"os"
	"path"
	"io"

	"net/http"
	"io/ioutil"

	"os/exec"

	"encoding/json"
	"gopkg.in/yaml.v2"
)

type InitTasks struct {
	log Log

	root string

	tasks[] InitTask
}

func (tasks *InitTasks) RunTasks() {
	for _, task := range tasks.tasks {
		tasks.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "INIT TASK:", task)
		task.RunTask(tasks.log)
	}
}

func (tasks *InitTasks) AddTask(task InitTask) {
	if tasks.tasks == nil {
		tasks.tasks = [] InitTask {}
	}

	tasks.tasks = append(tasks.tasks, task)
}

func (tasks *InitTasks) AddFile(path string, contents string) {
	tasks.AddTask(InitTask( & InitTaskFile {
		root: tasks.root,
		path: path,
		contents: contents,
	}))
}
func (tasks *InitTasks) AddRemoteFile(path string, url string) {
	tasks.AddTask(InitTask( & InitTaskRemoteFile {
		root: tasks.root,
		path: path,
		url: url,
	}))
}
func (tasks *InitTasks) AddFileCopy(path string, source string) {
	tasks.AddTask(InitTask( & InitTaskFileCopy {
		root: tasks.root,
		path: path,
		source: source,
	}))
}
func (tasks *InitTasks) AddGitClone(path string, url string) {
	tasks.AddTask(InitTask( & InitTaskGitClone {
		root: tasks.root,
		path: path,
		url: url,
	}))
}
func (tasks *InitTasks) AddMessage(message string) {
	tasks.AddTask(InitTask( & InitTaskMessage {
		message: message,
	}))
}
func (tasks *InitTasks) AddError(error string) {
	tasks.AddTask(InitTask( & InitTaskError {
		error: error,
	}))
}

type InitTask interface {
	RunTask(log Log) bool
}

type InitTaskFileBase struct {
	root string
}

func (task * InitTaskFileBase) MakeDir(log Log, makePath string, pathIsFile bool) bool {
	pd := path.Join(task.root, makePath)
	if pathIsFile {
		pd = path.Dir(pd)
	}

	if err := os.MkdirAll(pd, 0777);err != nil {
		// @todo something log
		return false
	}
	return true
}
func (task * InitTaskFileBase) MakeFile(log Log, destinationPath string, contents string) bool {
	if !task.MakeDir(log, destinationPath, true) {
		// @todo something log
		return false
	}

	pd := path.Join(task.root, destinationPath)

	fileObject, err := os.Create(pd)
	defer fileObject.Close()
	if err != nil {
		// @todo something log
		return false
	}
	if _, err := fileObject.WriteString(contents);
	err != nil {
		// @todo something log
		return false
	}

	return true
}

func (task * InitTaskFileBase) CopyFile(log Log, destinationPath string, sourcePath string) bool {

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		log.Warning("could not copy file as it does not exist [" + sourcePath + "] : " + err.Error())
		return false
	}
	defer sourceFile.Close()

	destinationRootPath := path.Join(task.root, destinationPath)
	if !task.MakeDir(log, destinationRootPath, true) {
		// @todo something log
		log.Warning("could not copy file as the path to the destination file could not be created [" + destinationPath + "]")
		return false
	}

	destinationFile, err := os.Open(destinationRootPath)
	if err == nil {
		log.Warning("could not copy file as it already exists [" + destinationPath + "]")
		defer destinationFile.Close()
		return false
	}

	destinationFile, err = os.Create(destinationRootPath)
	if err != nil {
		log.Warning("could not copy file as destination file could not be created [" + destinationPath + "] : " + err.Error())
		return false
	}

	defer destinationFile.Close()

	_, err = io.Copy(sourceFile, destinationFile)
	if err == nil {
		sourceInfo, err := os.Stat(sourcePath)
		if err == nil {
			err = os.Chmod(destinationPath, sourceInfo.Mode())
			return true
		} else {
			log.Warning("could not copy file as destination file could not be created [" + destinationPath + "] : " + err.Error())
			return false
		}
	} else {
		log.Warning("could not copy file as copy failed [" + destinationPath + "] : " + err.Error())
	}


	return true
}

func (task * InitTaskFileBase) CopyRemoteFile(log Log, destinationPath string, sourcePath string) bool {

	response, err := http.Get(sourcePath)
	if err != nil {
		log.Warning("Could not open remote URL: " + sourcePath)
		return false
	}
	defer response.Body.Close()

	sourceContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Warning("Could not read remote file: " + sourcePath)
		return false
	}

	return task.MakeFile(log, destinationPath, string(sourceContent))
}


func (task * InitTaskFileBase) CopyFileRecursive(log Log, path string, source string) bool {
	return task.copyFileRecursive(log, path, source, "")
}
func (task * InitTaskFileBase) copyFileRecursive(log Log, destinationRootPath string, sourceRootPath string, sourcePath string) bool {

	fullPath := sourceRootPath

		if sourcePath != "" {
		fullPath = path.Join(fullPath, sourcePath)
	}

	// get properties of source dir
	info,
	err := os.Stat(fullPath)
	if err != nil {
		// @TODO do something log : source doesn't exist
		return false
	}

	mode := info.Mode()
	if mode.IsDir() {

		directory, _ := os.Open(fullPath)
		objects, err := directory.Readdir(-1)

		if err != nil {
			// @TODO do something log : source doesn't exist
			return false
		}

		for _, obj := range objects {

			//childSourcePath := source + "/" + obj.Name()
			childSourcePath := path.Join(sourcePath, obj.Name())
			task.copyFileRecursive(log, destinationRootPath, sourceRootPath, childSourcePath)

		}

	} else {
		// add file copy
		destinationPath := path.Join(destinationRootPath, sourcePath)
		if task.CopyFile(log, destinationPath, sourceRootPath) {
			log.Info("--> Copied file (recursively): " + sourcePath + " [from " + sourceRootPath + "]")
			return true
		} else {
			log.Warning("--> Failed to copy file: " + sourcePath + " [from " + sourceRootPath + "]")
			return false
		}
		return true
	}
	return true
}


type InitTaskFile struct {
	InitTaskFileBase
	root string

	path string
	contents string
}
func (task * InitTaskFile) RunTask(log Log) bool {
	if task.MakeFile(log, task.path, task.contents) {
		log.Message("--> Created file : " + task.path)
		return true
	} else {
		log.Warning("--> Failed to create file : " + task.path)
		return false
	}
}

type InitTaskRemoteFile struct {
	InitTaskFileBase
	root string

	path string
	url string
}
func (task * InitTaskRemoteFile) RunTask(log Log) bool {
	if task.CopyRemoteFile(log, task.path, task.url) {
		log.Message("--> Copied remote file : " + task.url + " -> " + task.path)
		return true
	} else {
		log.Warning("--> Failed to copy remote file : " + task.url)
		return false
	}
}

type InitTaskFileCopy struct {
	InitTaskFileBase
	root string

	path string
	source string
}
func (task * InitTaskFileCopy) RunTask(log Log) bool {
	if task.CopyFileRecursive(log, task.path, task.source) {
		log.Message("--> Copied file : " + task.source + " -> " + task.path)
		return true
	} else {
		log.Warning("--> Failed to copy file : " + task.source + " -> " + task.path)
		return false
	}
}

type InitTaskGitClone struct {
	InitTaskFileBase
	root string

	path string
	url string
}
func (task * InitTaskGitClone) RunTask(log Log) bool {

	destinationPath := path.Join(task.root, task.path)
	url := task.url

	if !task.MakeDir(log, destinationPath, false) {
		return false
	}

	cmd := exec.Command("git", "clone", "--progress", url, destinationPath)
	cmd.Stderr = log
	err := cmd.Start()

	if err != nil {
		log.Error("Failed to clone the remote repository [" + url + "] => " + err.Error())
		return false
	}

	err = cmd.Wait()

	if err != nil {
		log.Error("Failed to clone the remote repository [" + url + "] => " + err.Error())
		return false
	}

	log.Message("Cloned remote repository [" + url + "] to local path " + destinationPath)
	return true
}

type InitTaskError struct {
	error string
}
func (task * InitTaskError) RunTask(log Log) bool {
	log.Error(task.error)
	return true
}

type InitTaskMessage struct {
	message string
}
func (task * InitTaskMessage) RunTask(log Log) bool {
	log.Message(task.message)
	return true
}


/**
 * Getting tasks from YAML
 */

func (tasks *InitTasks) AddTasksFromYaml(yamlSource[] byte) error {

	var yaml_tasks[] map[string] interface {}
	err := yaml.Unmarshal(yamlSource, & yaml_tasks)
	if err != nil {
		return err
	}

	var taskAdder TaskAdder
	for _, task_struct := range yaml_tasks {

		taskAdder = nil

		if _, ok := task_struct["Type"];
		!ok {
			continue
		}

		switch task_struct["Type"] {
			case "File":
				json_task, _ := json.Marshal(task_struct)
				var task InitTaskYaml_FileMake
				if err := json.Unmarshal(json_task, & task);
				err == nil {
					taskAdder = TaskAdder( & task)
				}
			case "RemoteFile":
				json_task, _ := json.Marshal(task_struct)
				var task InitTaskYaml_RemoteFileCopy
				if err := json.Unmarshal(json_task, & task);
				err == nil {
					taskAdder = TaskAdder( & task)
				}
			case "FileCopy":
				json_task, _ := json.Marshal(task_struct)
				var task InitTaskYaml_FileCopy
				if err := json.Unmarshal(json_task, & task);
				err == nil {
					taskAdder = TaskAdder( & task)
				}
			case "GitClone":
				json_task, _ := json.Marshal(task_struct)
				var task InitTaskYaml_GitClone
				if err := json.Unmarshal(json_task, & task);
				err == nil {
					taskAdder = TaskAdder( & task)
				}
			case "Message":
				json_task, _ := json.Marshal(task_struct)
				var task InitTaskYaml_Message
				if err := json.Unmarshal(json_task, & task);
				err == nil {
					taskAdder = TaskAdder( & task)
				}
			case "Error":
				json_task, _ := json.Marshal(task_struct)
				var task InitTaskYaml_Error
				if err := json.Unmarshal(json_task, & task);
				err == nil {
					taskAdder = TaskAdder( & task)
				}

			default:
				taskType := task_struct["Type"].(string)
				tasks.log.Warning("Unknown init task type [" + taskType + "]")
		}

		if taskAdder != nil {
			taskAdder.AddTask(tasks)
		}

	}

	return nil
}

type InitTaskYaml_Base struct {
	Type string `json:"Type" yaml:"Type"`
}

type TaskAdder interface {
	AddTask(tasks *InitTasks)
}

type InitTaskYaml_FileMake struct {
	Path string `json:"Path" yaml:"Path"`
	Contents string `json:"Contents" yaml:"Contents"`
}
func (task * InitTaskYaml_FileMake) AddTask(tasks *InitTasks) {
	tasks.AddFile(task.Path, task.Contents)
}

type InitTaskYaml_RemoteFileCopy struct {
	Path string `json:"Path" yaml:"Path"`
	Url string `json:"Url" yaml:"Url"`
}
func (task * InitTaskYaml_RemoteFileCopy) AddTask(tasks *InitTasks) {
	tasks.AddRemoteFile(task.Path, task.Url)
}

type InitTaskYaml_FileCopy struct {
	Path string `json:"Path" yaml:"Path"`
	Source string `json:"Source" yaml:"Source"`
}
func (task * InitTaskYaml_FileCopy) AddTask(tasks *InitTasks) {
	tasks.AddFileCopy(task.Path, task.Source)
}

type InitTaskYaml_GitClone struct {
	Path string `json:"Path" yaml:"Path"`
	Url string `json:"Url" yaml:"Url"`
}
func (task * InitTaskYaml_GitClone) AddTask(tasks *InitTasks) {
	tasks.AddGitClone(task.Path, task.Url)
}

type InitTaskYaml_Message struct {
	Message string `json:"Message" yaml:"Message"`
}
func (task * InitTaskYaml_Message) AddTask(tasks *InitTasks) {
	tasks.AddMessage(task.Message)
}

type InitTaskYaml_Error struct {
	Error string `json:"Error" yaml:"Error"`
}
func (task * InitTaskYaml_Error) AddTask(tasks *InitTasks) {
	tasks.AddError(task.Error)
}