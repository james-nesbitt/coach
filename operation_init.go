package main

import (
	"os"
	"path"
)

type Operation_Init struct {
	conf *Conf
	log Log

	root string
	variant string
	force bool

	Targets []string
}
func (operation *Operation_Init) Flags(flags []string) {

	operation.variant = "default"
	operation.force = false

	remainingFlags := []string{}
	for index, flag := range flags {
		switch flag {
			case "-f":
				fallthrough
			case "--force":
				operation.force = true

			default:
				remainingFlags = flags[index:]
				break;
		}
	}

	if variant := remainingFlags[0]; variant!="" {
		operation.variant = variant
	}

	operation.root, _ = os.Getwd()
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

	// get a list of files to create
	files := map[string]string{}
	switch operation.variant {
		case "default":
			fallthrough
		default:
			files = operation.Init_Default_Files()
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
