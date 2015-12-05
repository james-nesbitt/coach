package main

import (
	"path"
	"strings"

	"os"
	"os/exec"
 	"io/ioutil"

 	"net/http"
)

var (
	// coach demos are keyed remoteyamls
	COACH_DEMO_URLS = map[string]string{
		"lamp": "https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/lamp/.coach/coachinit.yml",
		"lamp_monolithic": "https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/lamp_monolithic/.coach/coachinit.yml",
		"lamp_multiplephps": "https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/lamp_multiplephps/.coach/coachinit.yml",
		"lamp_scaling": "https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/lamp_scaling/.coach/coachinit.yml",
	}
)

type Operation_Init struct {
	conf *Conf
	log Log

	root string						// Path to root

	handler strings 			// which method to use to initialize
	source string					// which method variant to initialize

	force bool

	targets []string			// often not used, but perhaps it might be useful to define what nodes to create in some scenarios
}

func (operation *Operation_Init) Flags(flags []string) {

	operation.root, _ = operation.conf.Path("projectroot")
	if operation.root=="" {
		operation.root, _ = os.Getwd()
	}

	operation.handler = "default"
	operation.source = ""

	remainingFlags := []string{}

	flagLoop:
		for index:=0; index<len(flags); index++ {
			flag:= flags[index]

			switch flag {
				case "-d":
					fallthrough
				case "--in-directory":
					if !strings.HasPrefix(flags[index+1], "-") {
						operation.root = flags[index+1]
						index++
					}
				case "-f":
					fallthrough
				case "--force":
					operation.force = true

				default:
					remainingFlags = flags[index:]
					break flagLoop
			}
		}

	if len(remainingFlags)>0 {
		handler := remainingFlags[0]
		remainingFlags = remainingFlags[1:]

		switch handler {
			case "demo":
				operation.handler = "remoteyaml"
				if len(remainingFlags)>0 {
					operation.source = remainingFlags[0]
					remainingFlags = remainingFlags[1:]
				} else {
					operation.source = "wunder"
				}
				operation.source = COACH_DEMO_URLS[operation.source]

			case "user":
				fallthrough
			case "yaml":
				fallthrough
			case "remoteyaml":
				fallthrough
			case "git":
				fallthrough
			case "default":

				operation.handler = handler
				if len(remainingFlags)>0 {
					operation.source = remainingFlags[0]
					remainingFlags = remainingFlags[1:]
				}

			default:
				operation.handler = "default"
				remainingFlags = append([]string{handler}, remainingFlags...)

				if len(remainingFlags)>0 {
					operation.source = remainingFlags[0]
					remainingFlags = remainingFlags[1:]
				}

		}

	}

}

func (operation *Operation_Init) Help(topics []string) {
	operation.log.Note(`Operation: INIT

Coach will attempt to initialize a new coach project in the current folder.

SYNTAX:
    $/> coach init [ {type} {type flags} ]

EXAMPLES:

    $/> coach init
    $/> coach init default
    Populate the current path with default settings

    $/> coach init user {template}

    Uses the contents of ~/.coach/templates to populate the current folder

    $/> coach init git https://github.com/aleksijohansson/docker-drupal-coach.git

    Clones the target git URL to the current path

    There are also various demo inits:
      $/> coach init demo lamp

      	creates a standard LAMP stack

      $/> coach init demo lamp_monolithic

        creates a single container LAMP stack

      $/> coach init demo lamp_multiplephps

        creates a LAMP stack with multiple php servers to

      $/> coach init demo lamp_scaling

YAML base inits

There is is a YAML file syntax that can be used to define an init,
which describes a set of operations such as creating fields etc. The 
demo inits are based on this forumla, and can be seen on the git repo
for the coach project.

If you have a fileset that you like, then convert it to the YAML syntax
and keep it on the internet, either in a repo, or a gist or even a pastebin.

Then you can create a new project using :
  
    $/> coach init remoteyaml http://path.to.my/yaml.yml

    (note that the path has to be a full body yml file)
`)
}

func (operation *Operation_Init) Run() {
	operation.log.Info("running init operation")

	var err error
	var ok bool
	var targetPath, coachPath string

	targetPath = operation.root
	if targetPath=="" {
		targetPath, ok = operation.conf.Path("project");
		if !ok || targetPath=="" {
			targetPath, err = os.Getwd()
			if err!=nil {
				operation.log.Error("No path suggested for new project init")
				return
			}
		}
	}

	_, err = os.Stat( targetPath )
	if err!=nil {
		operation.log.Error("Invalid path suggested for new project init : ["+targetPath+"] => "+err.Error())
		return
	}

	coachPath = path.Join(targetPath, coachConfigFolder);

	operation.log.Message("Preparing INIT operation ["+operation.handler+":"+operation.source+"] in path : "+targetPath)

	_, err = os.Stat( coachPath )
	if (!operation.force && err==nil) {
		operation.log.Error("cannot create new project folder, as one already exists")
		return
	}

	operation.log = operation.log.ChildLog(strings.ToUpper(operation.handler))
	tasks := InitTasks{log: operation.log.ChildLog("TASKS")}

	ok = true
	switch operation.handler {
		case "user":
			ok = operation.Init_User_Run(operation.source, &tasks)
		case "git":
			ok = operation.Init_Git_Run(operation.source, &tasks)
		case "yaml":
			ok = operation.Init_Yaml_Run(operation.source, &tasks)
		case "remoteyaml":
			ok = operation.Init_RemoteYaml_Run(operation.source, &tasks)
		case "default":
			ok = operation.Init_Default_Run(operation.source, &tasks)
		default:

			operation.log.Error("Unknown init handler "+operation.handler)
			ok = false

	}

	if ok {
		operation.log.Info("Running init tasks")
		tasks.RunTasks()
	} else {
		operation.log.Warning("No init tasks were defined.")
	}

}


func (operation *Operation_Init) Init_User_Run(template string, tasks *InitTasks) bool {

	if template=="" {
		operation.log.Error("You have not provided a template name  $/> coach init user {template}")
		return false
	}

	templatePath, ok := operation.conf.Path("usertemplates")
	if !ok {
		operation.log.Error("COACH has no user template path for the current user")
		return false
	}
	sourcePath := path.Join(templatePath , template )

	if _, err := os.Stat( sourcePath ); err!=nil {
		operation.log.Error("Invalid template path suggested for new project init : ["+template+"] expected path ["+sourcePath+"] => "+err.Error())
		return false
	}
	
	operation.log.Message("Perfoming init operation from user template ["+template+"] : "+sourcePath)

	tasks.AddFileCopy(operation.root, sourcePath)

	tasks.AddMessage("Copied coach template ["+template+"] to init project")
	tasks.AddFile(".coach/CREATEDFROM.md", `THIS PROJECT WAS CREATED FROM A User Template :`+template)

	return true
}

func (operation *Operation_Init) Init_Git_Run(source string, tasks *InitTasks) bool {

	if source=="" {
		operation.log.Error("You have not provided a git target $/> coach init git https://github.com/aleksijohansson/docker-drupal-coach")
		return false
	}

	url := source
	path := operation.root

	cmd := exec.Command("git", "clone", "--progress", url, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = operation.log
	cmd.Stderr = operation.log

	err := cmd.Start()

	if err!=nil {
		operation.log.Error("Failed to clone the remote repository ["+url+"] => "+err.Error())
		return false
	}

	operation.log.Message("Clone remote repository to local project folder ["+url+"]")
	err = cmd.Wait()

	if err!=nil {
		operation.log.Error("Failed to clone the remote repository ["+url+"] => "+err.Error())
		return false
	}

	tasks.AddMessage("Cloned remote repository ["+url+"] to local project folder")
	tasks.AddFile(".coach/CREATEDFROM.md", `THIS PROJECT WAS CREATED FROM GIT`)

	return true
}

// Get tasks from remote YAML corresponding to a remote yaml file
func (operation *Operation_Init) Init_Yaml_Run(path string, tasks *InitTasks) bool {

	// read the config file
	yamlSourceBytes, err := ioutil.ReadFile(path)
	if err!=nil {
		operation.log.Error("Could not read the YAML file ["+path+"]: "+err.Error())
		return false
	}
	if len(yamlSourceBytes)==0 {
		operation.log.Error("Yaml file ["+path+"] was empty")
		return false
	}

	tasks.AddMessage("Initializing using YAML Source ["+path+"] to local project folder")

	// get tasks from yaml
	tasks.AddTasksFromYaml(yamlSourceBytes)

	// Add some message items
	tasks.AddFile(".coach/CREATEDFROM.md", "THIS PROJECT WAS CREATED A COACH YAML INSTALLER :"+path)

	return true
}

// Get tasks from remote YAML corresponding to a remote yaml file
func (operation *Operation_Init) Init_RemoteYaml_Run(url string, tasks *InitTasks) bool {

	resp, err := http.Get(url)
	if err != nil {
		operation.log.Error("Could not retrieve remote yaml init instructions ["+url+"] : "+err.Error())
		return false
	}
	defer resp.Body.Close()
	yamlSourceBytes, err := ioutil.ReadAll(resp.Body)

	tasks.AddMessage("Initializing using Remote YAML Source ["+url+"] to local project folder")

	// get tasks from yaml
	tasks.AddTasksFromYaml(yamlSourceBytes)

	// Add some message items
	tasks.AddFile(".coach/CREATEDFROM.md", "THIS PROJECT WAS CREATED A COACH YAML INSTALLER :"+url)

	return true
}
