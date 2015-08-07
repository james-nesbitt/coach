package main

import (
 	git "github.com/libgit2/git2go"
)

func (operation *Operation_Init) Init_Git_Run(flags []string) (bool, map[string]string) {

	if len(flags)==0 {
		operation.log.Error("You have not provided a git target $/> coach init git https://github.com/aleksijohansson/docker-drupal-coach")
		return false, map[string]string{}
	}

	target := flags[0]

	options := git.CloneOptions{
		Bare: false,
	}

	if len(flags)>1 {
		options.CheckoutBranch = flags[1]
	}

	operation.log.Message("Clone remote repository to local project folder ["+target+"]")
	_, err := git.Clone(target, operation.root, &options)
	if err!=nil {
		operation.log.Error("Failed to clone the remote repository ["+target+"] => "+err.Error())
	} else {
		operation.log.Message("Clone remote repository to local project folder")
	}

	return true, map[string]string{
		".coach/CREATEDFROM.md":  `THIS PROJECT WAS CREATED FROM GIT`,
	}
}
