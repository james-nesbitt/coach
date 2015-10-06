package main

import (
	"os"
	"path"
	"strings"
)

type Operation_Init struct {
	conf *Conf
	log Log

	root string							// Path to root
	variant string					// which method to use to initialize
	handlerFlags []string 	// flags to pass to the variant handler

	force bool

	targets []string				// often not used, but perhaps it might be useful to define what nodes to create in some scenarios
}

func (operation *Operation_Init) Flags(flags []string) {

	operation.root, _ = os.Getwd()
	operation.variant = "default"

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
		if variant := remainingFlags[0]; variant!="" {
			operation.variant = variant
			remainingFlags = remainingFlags[1:]
		}
	}

	operation.handlerFlags = remainingFlags
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

`)
}

func (operation *Operation_Init) Run() {
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
	_, err = os.Stat( coachPath )
	if (!operation.force && err==nil) {
		operation.log.Error("cannot create new project folder, as one already exists")
		return
	}

	operation.log.Message("Preparing INIT operation in path : "+targetPath)

	// a list of files to create, optionally passed as a return by the init methods
	var files map[string]string

	operation.log = operation.log.ChildLog(strings.ToUpper(operation.variant))
	switch operation.variant {
		case "user":
			_, files = operation.Init_User_Run(operation.handlerFlags)
		case "git":
			_, files = operation.Init_Git_Run(operation.handlerFlags)
		case "default":
			fallthrough
		default:
			_, files = operation.Init_Default_Run(operation.handlerFlags)
	}

	for filePath, fileContents := range files {
		operation.log.Message("--> Create file : "+filePath)
		operation.MakeFile(filePath, fileContents)
	}
}


func (operation *Operation_Init) MakeFile(filePath string, fileContents string) bool {

	initPath := operation.root
	pd := path.Join(initPath, filePath)
	pdbase := path.Dir(pd)

	if err := os.MkdirAll(pdbase, 0777); err!=nil {
		return false
	}

	fileObject, err := os.Create(pd)
	defer fileObject.Close()
	if err!=nil {
		return false
	}
	if _, err := fileObject.WriteString(fileContents); err!=nil {
		return false
	}

	return true
}
