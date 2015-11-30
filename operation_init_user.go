package main

import (
	"os"
	"path"
)

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
