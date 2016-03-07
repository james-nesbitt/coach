package conf

import (
	"strconv"

	"github.com/james-nesbitt/coach/log"
)

const (
	COACH_CONF_ENVIRONMENTS_DEFAULT = "default"
	COACH_CONF_ENVIRONMENTS_SUBPATH = "environments"
)

var (
	environmentIncrement = 0
)

// Look for project configurations inside the environment conf paths
func (project *Project) from_EnvironmentsPath(logger log.Log) {
	for _, yamlEnvironmentPath := range project.Paths.GetConfSubPaths(COACH_CONF_ENVIRONMENTS_SUBPATH) {
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Looking for Environment subpath: "+yamlEnvironmentPath)

		if project.CheckFileExists(yamlEnvironmentPath) {

			logger.Debug(log.VERBOSITY_DEBUG, "ADDING Environment subpath: "+yamlEnvironmentPath)

			environmentIncrement++
			envid := COACH_CONF_ENVIRONMENTS_SUBPATH+"-"+strconv.Itoa(environmentIncrement)

			project.SetPath(envid, yamlEnvironmentPath, true)
			project.setConfPath("project-coach")

		}
	}
}
