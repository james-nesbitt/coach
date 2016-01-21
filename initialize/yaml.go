package initialize

import (
	"io/ioutil"
	"net/http"
	"strings"

	"encoding/json"
	"gopkg.in/yaml.v2"

	"github.com/james-nesbitt/coach-tools/log"
)

// Get tasks from remote YAML corresponding to a remote yaml file
func (tasks *InitTasks) Init_Yaml_Run(logger log.Log, path string) bool {

	var yamlSourceBytes []byte
	var err error

	if strings.Contains(path, "://") {

		resp, err := http.Get(path)
		if err != nil {
			logger.Error("Could not retrieve remote yaml init instructions [" + path + "] : " + err.Error())
			return false
		}
		defer resp.Body.Close()
		yamlSourceBytes, err = ioutil.ReadAll(resp.Body)

	} else {

		// read the config file
		yamlSourceBytes, err = ioutil.ReadFile(path)
		if err != nil {
			logger.Error("Could not read the local YAML file [" + path + "]: " + err.Error())
			return false
		}
		if len(yamlSourceBytes) == 0 {
			logger.Error("Yaml file [" + path + "] was empty")
			return false
		}

	}

	tasks.AddMessage("Initializing using YAML Source [" + path + "] to local project folder")

	// get tasks from yaml
	tasks.AddTasksFromYaml(logger, yamlSourceBytes)

	// Add some message items
	tasks.AddFile(".coach/CREATEDFROM.md", "THIS PROJECT WAS CREATED A COACH YAML INSTALLER :"+path)

	return true
}

/**
 *Getting tasks from YAML
 */

func (tasks *InitTasks) AddTasksFromYaml(logger log.Log, yamlSource []byte) error {

	var yaml_tasks []map[string]interface{}
	err := yaml.Unmarshal(yamlSource, &yaml_tasks)
	if err != nil {
		return err
	}

	var taskAdder TaskAdder
	for _, task_struct := range yaml_tasks {

		taskAdder = nil

		if _, ok := task_struct["Type"]; !ok {
			continue
		}

		switch task_struct["Type"] {
		case "File":
			json_task, _ := json.Marshal(task_struct)
			var task InitTaskYaml_FileMake
			if err := json.Unmarshal(json_task, &task); err == nil {
				taskAdder = TaskAdder(&task)
			}
		case "RemoteFile":
			json_task, _ := json.Marshal(task_struct)
			var task InitTaskYaml_RemoteFileCopy
			if err := json.Unmarshal(json_task, &task); err == nil {
				taskAdder = TaskAdder(&task)
			}
		case "FileCopy":
			json_task, _ := json.Marshal(task_struct)
			var task InitTaskYaml_FileCopy
			if err := json.Unmarshal(json_task, &task); err == nil {
				taskAdder = TaskAdder(&task)
			}
		case "GitClone":
			json_task, _ := json.Marshal(task_struct)
			var task InitTaskYaml_GitClone
			if err := json.Unmarshal(json_task, &task); err == nil {
				taskAdder = TaskAdder(&task)
			}
		case "Message":
			json_task, _ := json.Marshal(task_struct)
			var task InitTaskYaml_Message
			if err := json.Unmarshal(json_task, &task); err == nil {
				taskAdder = TaskAdder(&task)
			}
		case "Error":
			json_task, _ := json.Marshal(task_struct)
			var task InitTaskYaml_Error
			if err := json.Unmarshal(json_task, &task); err == nil {
				taskAdder = TaskAdder(&task)
			}

		default:
			taskType := task_struct["Type"].(string)
			logger.Warning("Unknown init task type [" + taskType + "]")
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
	Path     string `json:"Path" yaml:"Path"`
	Contents string `json:"Contents" yaml:"Contents"`
}

func (task *InitTaskYaml_FileMake) AddTask(tasks *InitTasks) {
	tasks.AddFile(task.Path, task.Contents)
}

type InitTaskYaml_RemoteFileCopy struct {
	Path string `json:"Path" yaml:"Path"`
	Url  string `json:"Url" yaml:"Url"`
}

func (task *InitTaskYaml_RemoteFileCopy) AddTask(tasks *InitTasks) {
	tasks.AddRemoteFile(task.Path, task.Url)
}

type InitTaskYaml_FileCopy struct {
	Path   string `json:"Path" yaml:"Path"`
	Source string `json:"Source" yaml:"Source"`
}

func (task *InitTaskYaml_FileCopy) AddTask(tasks *InitTasks) {
	tasks.AddFileCopy(task.Path, task.Source)
}

type InitTaskYaml_GitClone struct {
	Path string `json:"Path" yaml:"Path"`
	Url  string `json:"Url" yaml:"Url"`
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
