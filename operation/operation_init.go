package operation

import (
	"strings"

	"os"

	"github.com/james-nesbitt/coach-tools/conf"
	"github.com/james-nesbitt/coach-tools/initialize"
	"github.com/james-nesbitt/coach-tools/log"
)

type InitOperation struct {
	log  log.Log
	conf *conf.Project

	force bool

	root string // Path to root

	handler string // which method to use to initialize
	source  string // which method variant to initialize
}

func (operation *InitOperation) Id() string {
	return "init"
}
func (operation *InitOperation) Flags(flags []string) bool {
	operation.root, _ = operation.conf.Path("project-root")
	if operation.root == "" {
		operation.root, _ = os.Getwd()
	}

	operation.handler = "default"
	operation.source = ""

	remainingFlags := []string{}

flagLoop:
	for index := 0; index < len(flags); index++ {
		flag := flags[index]

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

	if len(remainingFlags) > 0 {
		handler := remainingFlags[0]
		remainingFlags = remainingFlags[1:]

		switch handler {
		case "demo":
			fallthrough
		case "yaml":
			fallthrough
		case "remoteyaml":
			fallthrough
		case "user":
			fallthrough
		case "git":
			fallthrough
		case "default":

			operation.handler = handler
			if len(remainingFlags) > 0 {
				operation.source = remainingFlags[0]
				remainingFlags = remainingFlags[1:]
			}

		default:
			operation.handler = "default"
			remainingFlags = append([]string{handler}, remainingFlags...)

			if len(remainingFlags) > 0 {
				operation.source = remainingFlags[0]
				remainingFlags = remainingFlags[1:]
			}

		}

	}
	return true
}
func (operation *InitOperation) Help(topics []string) {
	operation.log.Message(`Operation: Init

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

		$/> coach init demo complete

		creates a really full and noisy example of a project

YAML base inits

There is is a YAML file syntax that can be used to define an init,
which describes a set of operations such as creating fields etc. The 
demo inits are based on this forumla, and can be seen on the git repo
for the coach project.

If you have a fileset that you like, then convert it to the YAML syntax
and keep it on the internet, either in a repo, or a gist or even a pastebin.

Then you can create a new project using :
  
	$/> coach init yaml http://path.to.my/yaml.yml

	(note that the path has to be a full body yml file)
`)
}

func (operation *InitOperation) Run(logger log.Log) bool {
	logger.Info("running init operation")

	var err error
	var ok bool
	var targetPath, coachPath string

	targetPath = operation.root
	if targetPath == "" {
		targetPath, ok = operation.conf.Path("project-root")
		if !ok || targetPath == "" {
			targetPath, err = os.Getwd()
			if err != nil {
				logger.Error("No path suggested for new project init")
				return false
			}
		}
	}

	_, err = os.Stat(targetPath)
	if err != nil {
		logger.Error("Invalid path suggested for new project init : [" + targetPath + "] => " + err.Error())
		return false
	}

	coachPath, _ = operation.conf.Paths.Path("coach-root")

	logger.Message("Preparing INIT operation [" + operation.handler + ":" + operation.source + "] in path : " + targetPath)

	_, err = os.Stat(coachPath)
	if !operation.force && err == nil {
		logger.Error("cannot create new project folder, as one already exists")
		return false
	}

	logger = logger.MakeChild(strings.ToUpper(operation.handler))
	tasks := initialize.InitTasks{}
	tasks.Init(logger.MakeChild("TASKS"), operation.conf, targetPath)

	ok = true
	switch operation.handler {
	case "user":
		ok = tasks.Init_User_Run(logger, operation.source)
	case "demo":
		ok = tasks.Init_Demo_Run(logger, operation.source)
	case "git":
		ok = tasks.Init_Git_Run(logger, operation.source)
	case "yaml":
		ok = tasks.Init_Yaml_Run(logger, operation.source)
	case "default":
		ok = tasks.Init_Default_Run(logger, operation.source)
	default:

		logger.Error("Unknown init handler " + operation.handler)
		ok = false

	}

	if ok {
		logger.Info("Running init tasks")
		tasks.RunTasks(logger)
		return true
	} else {
		logger.Warning("No init tasks were defined.")
		return false
	}

}
