package main

import (
	"os"
 	"os/exec"
)

func (operation *Operation_Init) Init_Git_Run(flags []string) (bool, map[string]string) {

	if len(flags)==0 {
		operation.log.Error("You have not provided a git target $/> coach init git https://github.com/aleksijohansson/docker-drupal-coach")
		return false, map[string]string{}
	}

	url := flags[0]
	path := operation.root

	cmd := exec.Command("git", "clone", "--progress", url, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = operation.log
	cmd.Stderr = operation.log


	err := cmd.Start()

	if err!=nil {
		operation.log.Error("Failed to clone the remote repository ["+url+"] => "+err.Error())
		return false, map[string]string{}
	}

	operation.log.Message("Clone remote repository to local project folder ["+url+"]")
	err = cmd.Wait()

	if err!=nil {
		operation.log.Error("Failed to clone the remote repository ["+url+"] => "+err.Error())
		return false, map[string]string{}
	} else {
		operation.log.Message("Cloned remote repository to local project folder")
		return true, map[string]string{
			".coach/CREATEDFROM.md":  `THIS PROJECT WAS CREATED FROM GIT`,
		}
	}

}
