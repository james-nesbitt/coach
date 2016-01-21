package initialize

import (
	"github.com/james-nesbitt/coach/log"
)

var (
	// coach demos are keyed remoteyamls
	COACH_DEMO_URLS = map[string]string{
		"complete":          "https://raw.githubusercontent.com/james-nesbitt/coach/master/init/templates/demo/complete/.coach/coachinit.yml",
		"drupal8":           "https://raw.githubusercontent.com/james-nesbitt/coach/master/init/templates/demo/drupal8/.coach/coachinit.yml",
		"lamp":              "https://raw.githubusercontent.com/james-nesbitt/coach/master/init/templates/demo/lamp/.coach/coachinit.yml",
		"lamp_monolithic":   "https://raw.githubusercontent.com/james-nesbitt/coach/master/init/templates/demo/lamp_monolithic/.coach/coachinit.yml",
		"lamp_multiplephps": "https://raw.githubusercontent.com/james-nesbitt/coach/master/init/templates/demo/lamp_multiplephps/.coach/coachinit.yml",
		"lamp_scaling":      "https://raw.githubusercontent.com/james-nesbitt/coach/master/init/templates/demo/lamp_scaling/.coach/coachinit.yml",
	}
)

func (tasks *InitTasks) Init_Demo_Run(logger log.Log, demo string) bool {
	if demoPath, ok := COACH_DEMO_URLS[demo]; ok {
		return tasks.Init_Yaml_Run(logger, demoPath)
	} else {
		logger.Error("Unknown demo key : " + demo)
		return false
	}
}
