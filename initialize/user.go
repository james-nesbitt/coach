package initialize

import (
	"os"
	"path"

	"github.com/james-nesbitt/coach/log"
)

func (tasks *InitTasks) Init_User_Run(logger log.Log, template string) bool {

	if template == "" {
		logger.Error("You have not provided a template name  $/> coach init user {template}")
		return false
	}

	templatePath, ok := tasks.conf.Path("user-templates")
	if !ok {
		logger.Error("COACH has no user template path for the current user")
		return false
	}
	sourcePath := path.Join(templatePath, template)

	if _, err := os.Stat(sourcePath); err != nil {
		logger.Error("Invalid template path suggested for new project init : [" + template + "] expected path [" + sourcePath + "] => " + err.Error())
		return false
	}

	logger.Message("Perfoming init operation from user template [" + template + "] : " + sourcePath)

	tasks.AddFileCopy(tasks.root, sourcePath)

	tasks.AddMessage("Copied coach template [" + template + "] to init project")
	tasks.AddFile(".coach/CREATEDFROM.md", `THIS PROJECT WAS CREATED FROM A User Template :`+template)

	return true
}
