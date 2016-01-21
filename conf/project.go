package conf

import (
	"os"

	"github.com/james-nesbitt/coach-tools/log"
)

// MakeCoachProject Project constructor for building a project based on a project path
func MakeCoachProject(logger log.Log, workingDir string) (project *Project) {
	project = &Project{
		Paths:        MakePaths(),        // empty typesafe paths object
		Tokens:       MakeTokens(),       // empty tokens object
		ProjectFlags: MakeProjectFlags(), // empty flags object
	}

	/**
	 * 1. try to get some base default configuration, seeing if there is a
	 *    project .coach folder above the workingdir (or in it) or maybe
	 *		there's a user's ~/.coach
	 */
	project.from_DefaultPaths(logger.MakeChild("defaults"), workingDir)

	// set the two project paths, to be the most important
	project.setConfPath("project-coach")
	// set the user conf paths, to be the least important
	project.setConfPath("user-coach")

	/**
	 * 2. Look for YAML configurations in the Configuration Paths
	 */
	project.from_ConfYaml(logger.MakeChild("confyaml"))

	/**
	 * 3. Try to load secrets from configuration paths
	 */
	project.from_SecretsYaml(logger.MakeChild("secrets"))

	return
}

// Project settings handler for coach, used to centralize and validate settings for a project
type Project struct {
	Name   string
	Author string

	Paths

	Tokens

	ProjectFlags
}

// Is this project configured enough to run coach
func (project *Project) IsValid(logger log.Log) bool {
	/**
	 * 1. do we have a project coach folder, and does it exist.
	 */
	if projectCoachPath, ok := project.Path("project-coach"); ok {
		if _, err := os.Stat(projectCoachPath); err != nil {
			logger.Warning(`Could not find a project root .coach folder:  
- At the root of any coach prpject, must be a .coach folder;
- This folder can container project configurations;
- The folder is required, because it tells coach where the project base is.`)
			return false
		}
	} else {
		return false
	}

	/**
	 * 2. Do we have a project name
	 *
	 * This is important as it gets used to make image and container names
	 */
	if project.Name == "" {
		logger.Warning(`Coach project has no Name.  
- A project name can be set in the .coach/conf.yml file.  
- The Name is used as a base for image and container names.`)
		return false
	}

	return true
}
