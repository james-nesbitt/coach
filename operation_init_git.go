package main

import (
	"os"
 	"os/exec"
)

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
