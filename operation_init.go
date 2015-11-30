package main

import (
	"os"
	"path"
	"io"
	"strings"
)

type Operation_Init struct {
	conf *Conf
	log Log

	root string							// Path to root

  handler string          // which method to use to initialize
	source string					// which method variant to initialize

	force bool

	targets []string				// often not used, but perhaps it might be useful to define what nodes to create in some scenarios
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
			case "user":
				fallthrough
			case "git":
				fallthrough
			case "demo":
				fallthrough
			case "default":

				operation.handler = handler

			default:
				operation.handler = "default"
				remainingFlags = append([]string{handler}, remainingFlags...)

		}

		if len(remainingFlags)>0 {
			operation.source = remainingFlags[0]
			remainingFlags = remainingFlags[1:]
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
      $/> coach init demo scale
      $/> coach init demo wunder

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

	operation.log.Message("Preparing INIT operation ["+operation.handler+"/"+operation.source+"] in path : "+targetPath)

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
		case "demo":
			ok = operation.Init_Demo_Run(operation.source, &tasks)
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


type InitTasks struct {
  log Log

  root string

  tasks []InitTask
}

func (tasks *InitTasks) RunTasks() {
  for _, task := range tasks.tasks {
  	tasks.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "INIT TASK:", task)
  	task.RunTask(tasks.log)
  }	
}

func (tasks *InitTasks) AddTask(task InitTask) {
	if tasks.tasks==nil {
		tasks.tasks = []InitTask{}
	}

	tasks.tasks = append(tasks.tasks, task)
}

func (tasks *InitTasks) AddFile(path string, contents string) {
	tasks.AddTask( InitTask( &InitTaskFile{
		root: tasks.root,
		path: path,
		contents: contents,
	} ))	
}
func (tasks *InitTasks) AddFileCopy(path string, source string) {
	tasks.AddTask( InitTask( &InitTaskFileCopy{
		root: tasks.root,
		path: path,
		source: source,
	} ))		
}
func (tasks *InitTasks) AddMessage(message string) {
	tasks.AddTask( InitTask( &InitTaskMessage{
		message: message,
	} ))
}
func (tasks *InitTasks) AddError(error string) {
	tasks.AddTask( InitTask( &InitTaskError{
		error: error,
	} ))
}

type InitTask interface {
	RunTask(log Log) bool
}

type InitTaskFileBase struct {
	root string
}

func (task *InitTaskFileBase) MakeDir(log Log, makePath string, pathIsFile bool) bool {
	pd := path.Join(task.root, makePath)
	if pathIsFile {
		pd = path.Dir(pd)
	}

	if err := os.MkdirAll(pd, 0777); err!=nil {
		// @todo something log
		return false
	}	
	return true
}
func (task *InitTaskFileBase) MakeFile(log Log, makePath string, contents string) bool {
	if !task.MakeDir(log, makePath, true) {
    // @todo something log
		return false
	}

	pd := path.Join(task.root, makePath)

	fileObject, err := os.Create(pd)
	defer fileObject.Close()
	if err!=nil {
    // @todo something log
		return false
	}
	if _, err := fileObject.WriteString(contents); err!=nil {
	  // @todo something log
		return false
	}

	return true
}

func (task *InitTaskFileBase) CopyFile(log Log, destinationPath string, sourcePath string) bool {

	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		log.Warning("could not copy file as it does not exist ["+sourcePath+"] : "+err.Error())
		return false
	}	
	defer sourceFile.Close()

	destinationRootPath := path.Join(task.root, destinationPath)
	if !task.MakeDir(log, destinationRootPath, true) {
    // @todo something log
		log.Warning("could not copy file as the path to the destination file could not be created ["+destinationPath+"]")
		return false
	}

	destinationFile, err := os.Open(destinationRootPath)
	if err == nil {
		log.Warning("could not copy file as it already exists ["+destinationPath+"]")
		defer destinationFile.Close()
		return false
	}

	destinationFile, err = os.Create(destinationRootPath)
	if err != nil {
		log.Warning("could not copy file as destination file could not be created ["+destinationPath+"] : "+err.Error())
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
			log.Warning("could not copy file as destination file could not be created ["+destinationPath+"] : "+err.Error())
			return false
		}
	} else {
		log.Warning("could not copy file as copy failed ["+destinationPath+"] : "+err.Error())
	}


	return true
}

func (task *InitTaskFileBase) CopyFileRecursive(log Log, path string, source string) bool {
	return task.copyFileRecursive(log, path, source, "")
}
func (task *InitTaskFileBase) copyFileRecursive(log Log, destinationRootPath string, sourceRootPath string, sourcePath string) bool {

		fullPath := sourceRootPath

		if sourcePath!="" {
			fullPath = path.Join(fullPath, sourcePath)
		}

		// get properties of source dir
		info, err := os.Stat(fullPath)
		if err!=nil {
			// @TODO do something log : source doesn't exist
			return false
		}

		mode := info.Mode()
		if mode.IsDir() {

			directory, _ := os.Open(fullPath)
			objects, err := directory.Readdir(-1)

			if err!=nil {
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
		  if task.CopyFile(log, destinationPath, sourceRootPath ) {
		  	log.Message("--> Copied file (recursively): "+sourcePath+" [from "+sourceRootPath+"]")
		  	return true
		  } else {
		  	log.Warning("--> Failed to copy file: "+sourcePath+" [from "+sourceRootPath+"]")
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
	if task.MakeFile(log, task.path, task.contents) {
		log.Message("--> Created file : "+task.path)
		return true
	} else {
		log.Warning("--> Failed to create file : "+task.path)
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
	if task.CopyFileRecursive(log, task.path, task.source) {
		log.Message("--> Copied file : "+task.source+" -> "+task.path)
		return true
	} else {
		log.Warning("--> Failed to copy file : "+task.source+" -> "+task.path)
		return false
	}
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
