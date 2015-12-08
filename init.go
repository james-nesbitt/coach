package main

import (
	"os"
	"os/user"
	"path"
	"io"
  "strings"

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
	tasks.AddTask(InitTask(&InitTaskFile {
		root: tasks.root,
		path: path,
		contents: contents,
	}))
}
func (tasks *InitTasks) AddRemoteFile(path string, url string) {
	tasks.AddTask(InitTask(&InitTaskRemoteFile {
		root: tasks.root,
		path: path,
		url: url,
	}))
}
func (tasks *InitTasks) AddFileCopy(path string, source string) {
	tasks.AddTask(InitTask(&InitTaskFileCopy {
		root: tasks.root,
		path: path,
		source: source,
	}))
}
func (tasks *InitTasks) AddGitClone(path string, url string) {
	tasks.AddTask(InitTask(&InitTaskGitClone {
		root: tasks.root,
		path: path,
		url: url,
	}))
}
func (tasks *InitTasks) AddMessage(message string) {
	tasks.AddTask(InitTask(&InitTaskMessage {
		message: message,
	}))
}
func (tasks *InitTasks) AddError(error string) {
	tasks.AddTask(InitTask(&InitTaskError {
		error: error,
	}))
}

type InitTask interface {
	RunTask(log Log) bool
}

type InitTaskFileBase struct {
	root string
}

func (task *InitTaskFileBase) absolutePath(targetPath string, addRoot bool) (string, bool) {
	if strings.HasPrefix(targetPath, "~") {
		return path.Join(task.userHomePath(), targetPath[1:]), !addRoot

	// I am not sure how reliable this function is
	// you passed an absolute path, so I can't add the root
	// } else if path.isAbs(targetPath) {
	// 	return targetPath, !addRoot

	// you passed a relative path, and want me to add the root
	} else if addRoot {
		return path.Join(task.root, targetPath), true

	// you passed path and don't want the root added (but is it already abs?)
	} else if targetPath!="" {
		return targetPath, true

	// you passed an empty string, and don't want the root added?
	} else {
		return targetPath, false
	}
}
func (task *InitTaskFileBase) userHomePath() (string) {
	if currentUser,  err := user.Current(); err==nil {
		return currentUser.HomeDir
	} else {
		return os.Getenv("HOME")
	}
}

func (task *InitTaskFileBase) MakeDir(log Log, makePath string, pathIsFile bool) bool {	
	if makePath=="" {
		return true // it's already made
	}

	if pathDirectory, ok := task.absolutePath(makePath, true); !ok {
		log.Warning("Invalid directory path: " + pathDirectory)
		return false
	}
	pathDirectory := path.Join(task.root, makePath)
	if pathIsFile {
		pathDirectory = path.Dir(pathDirectory)
	}

	if err := os.MkdirAll(pathDirectory, 0777);err != nil {
		// @todo something log
		return false
	}
	return true
}
func (task *InitTaskFileBase) MakeFile(log Log, destinationPath string, contents string) bool {
	if !task.MakeDir(log, destinationPath, true) {
		// @todo something log
		return false
	}

	if destinationPath, ok := task.absolutePath(destinationPath, true); !ok {
		log.Warning("Invalid file destination path: " + destinationPath)
		return false
	}

	fileObject, err := os.Create(destinationPath)
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

func (task *InitTaskFileBase) CopyFile(log Log, destinationPath string, sourcePath string) bool {
	if destinationPath=="" || sourcePath=="" {
		log.Warning("empty source or destination passed for copy")
		return false
	}

	sourcePath, ok := task.absolutePath(sourcePath, false)
	if !ok {
		log.Warning("Invalid copy source path: " + destinationPath)
		return false
	}
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		log.Warning("could not copy file as it does not exist [" + sourcePath + "] : " + err.Error())
		return false
	}
	defer sourceFile.Close()

	if !task.MakeDir(log, destinationPath, true) {
		// @todo something log
		log.Warning("could not copy file as the path to the destination file could not be created [" + destinationPath + "]")
		return false
	}
	destinationAbsPath, ok := task.absolutePath(destinationPath, true)
	if !ok {
		log.Warning("Invalid copy destination path: " + destinationPath)
		return false
	}

	destinationFile, err := os.Open(destinationAbsPath)
	if err == nil {
		log.Warning("could not copy file as it already exists [" + destinationPath + "]")
		destinationFile.Close()
		return false
	}

	destinationFile, err = os.Create(destinationAbsPath)
	if err != nil {
		log.Warning("could not copy file as destination file could not be created [" + destinationPath + "] : " + err.Error())
		return false
	}

	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)

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

func (task *InitTaskFileBase) CopyRemoteFile(log Log, destinationPath string, sourcePath string) bool {
	if destinationPath=="" || sourcePath=="" {
		return false
	}

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


func (task *InitTaskFileBase) CopyFileRecursive(log Log, path string, source string) bool {
	sourceAbsPath, ok := task.absolutePath(source, false)
	if !ok {
		log.Warning("Couldn't find copy source "+source)
		return false
	}
	return task.copyFileRecursive(log, path, sourceAbsPath, "")
}
func (task *InitTaskFileBase) copyFileRecursive(log Log, destinationRootPath string, sourceRootPath string, sourcePath string) bool {
	fullPath := sourceRootPath

	if sourcePath != "" {
		fullPath = path.Join(fullPath, sourcePath)
	}

	// get properties of source dir
	info,
	err := os.Stat(fullPath)
	if err != nil {
		// @TODO do something log : source doesn't exist
		log.Warning("File does not exist :"+fullPath)
		return false
	}

	mode := info.Mode()
	if mode.IsDir() {

		directory, _ := os.Open(fullPath)
		objects, err := directory.Readdir(-1)

		if err != nil {
			// @TODO do something log : source doesn't exist
			log.Warning("Could not open directory")
			return false
		}

		for _, obj := range objects {

			//childSourcePath := source + "/" + obj.Name()
			childSourcePath := path.Join(sourcePath, obj.Name())
			if !task.copyFileRecursive(log, destinationRootPath, sourceRootPath, childSourcePath) {
  			log.Warning("Resursive copy failed")
			}

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
func (task *InitTaskFile) RunTask(log Log) bool {
	if task.path=="" {
		return false
	}

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
func (task *InitTaskRemoteFile) RunTask(log Log) bool {
	if task.path=="" || task.root=="" || task.url=="" {
		return false
	}

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
func (task *InitTaskFileCopy) RunTask(log Log) bool {
	if task.path=="" || task.root=="" || task.source=="" {
		return false
	}

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
func (task *InitTaskGitClone) RunTask(log Log) bool {
	if task.root=="" || task.url=="" {
log.Error("EMPTY ROOT PASSED TO GIT: "+task.root)
		return false
	}

	destinationPath := task.path
	url := task.url

	if !task.MakeDir(log, destinationPath, false) {
		return false
	}

	destinationAbsPath, ok := task.absolutePath(destinationPath , true)
	if !ok {
		log.Warning("Invalid copy destination path: " + destinationPath)
		return false
	}

	cmd := exec.Command("git", "clone", "--progress", url, destinationAbsPath)
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
func (task *InitTaskError) RunTask(log Log) bool {
	log.Error(task.error)
	return true
}

type InitTaskMessage struct {
	message string
}
func (task *InitTaskMessage) RunTask(log Log) bool {
	log.Message(task.message)
	return true
}


/**
 *Getting tasks from YAML
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
					taskAdder = TaskAdder(&task)
				}
			case "RemoteFile":
				json_task, _ := json.Marshal(task_struct)
				var task InitTaskYaml_RemoteFileCopy
				if err := json.Unmarshal(json_task, & task);
				err == nil {
					taskAdder = TaskAdder(&task)
				}
			case "FileCopy":
				json_task, _ := json.Marshal(task_struct)
				var task InitTaskYaml_FileCopy
				if err := json.Unmarshal(json_task, & task);
				err == nil {
					taskAdder = TaskAdder(&task)
				}
			case "GitClone":
				json_task, _ := json.Marshal(task_struct)
				var task InitTaskYaml_GitClone
				if err := json.Unmarshal(json_task, & task);
				err == nil {
					taskAdder = TaskAdder(&task)
				}
			case "Message":
				json_task, _ := json.Marshal(task_struct)
				var task InitTaskYaml_Message
				if err := json.Unmarshal(json_task, & task);
				err == nil {
					taskAdder = TaskAdder(&task)
				}
			case "Error":
				json_task, _ := json.Marshal(task_struct)
				var task InitTaskYaml_Error
				if err := json.Unmarshal(json_task, & task);
				err == nil {
					taskAdder = TaskAdder(&task)
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
func (task *InitTaskYaml_FileMake) AddTask(tasks *InitTasks) {
	tasks.AddFile(task.Path, task.Contents)
}

type InitTaskYaml_RemoteFileCopy struct {
	Path string `json:"Path" yaml:"Path"`
	Url string `json:"Url" yaml:"Url"`
}
func (task *InitTaskYaml_RemoteFileCopy) AddTask(tasks *InitTasks) {
	tasks.AddRemoteFile(task.Path, task.Url)
}

type InitTaskYaml_FileCopy struct {
	Path string `json:"Path" yaml:"Path"`
	Source string `json:"Source" yaml:"Source"`
}
func (task *InitTaskYaml_FileCopy) AddTask(tasks *InitTasks) {
	tasks.AddFileCopy(task.Path, task.Source)
}

type InitTaskYaml_GitClone struct {
	Path string `json:"Path" yaml:"Path"`
	Url string `json:"Url" yaml:"Url"`
}
func (task *InitTaskYaml_GitClone) AddTask(tasks *InitTasks) {
	tasks.AddGitClone(task.Path, task.Url)
}

type InitTaskYaml_Message struct {
	Message string `json:"Message" yaml:"Message"`
}
func (task *InitTaskYaml_Message) AddTask(tasks *InitTasks) {
	tasks.AddMessage(task.Message)
}

type InitTaskYaml_Error struct {
	Error string `json:"Error" yaml:"Error"`
}
func (task *InitTaskYaml_Error) AddTask(tasks *InitTasks) {
	tasks.AddError(task.Error)
}