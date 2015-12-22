package conf

import (
	"os"
	"os/user"
	"path"

	"github.com/james-nesbitt/coach-tools/log"
)

const (
	COACH_PROJECT_CONF_FOLDER = ".coach"
	COACH_USER_CONF_SUBPATH = ".coach"  // @TODO move this to ~/.config/shared/coach ?
)

// Get some base configuration for the project conf based on
// searching for a .coach path
func (project *Project) from_DefaultPaths(logger log.Log, workingDir string) {

	logger.Debug(log.VERBOSITY_DEBUG_LOTS,"Creating default Project")

	homeDir := "."
	if currentUser,  err := user.Current(); err==nil {
		homeDir = currentUser.HomeDir
	} else {
		homeDir = os.Getenv("HOME")
	}

	projectRootDirectory := workingDir
	_, err := os.Stat( path.Join(projectRootDirectory, COACH_PROJECT_CONF_FOLDER) )
	RootSearch:
		for err!=nil {
			projectRootDirectory = path.Dir(projectRootDirectory)
			if (projectRootDirectory==homeDir || projectRootDirectory=="." || projectRootDirectory=="/") {
				logger.Info("Could not find a project folder, coach will assume that this project is not initialized.")
				projectRootDirectory = workingDir
				break RootSearch
			}
			_, err = os.Stat(path.Join(projectRootDirectory, COACH_PROJECT_CONF_FOLDER) )
		}

	/**
	 * Set up some frequesntly used paths
	 */
	project.SetPath("user-home", homeDir, true)
	project.SetPath("user-coach", path.Join(homeDir,COACH_PROJECT_CONF_FOLDER), true)
	project.SetPath("project-root", projectRootDirectory, true)
	project.SetPath("project-coach", path.Join(projectRootDirectory,COACH_PROJECT_CONF_FOLDER), true)

	/**
	 * @Note that it is advisable to not test if a path exists, as 
	 * that can cause race conditions, and can produce an invalid test
	 * as the path could be created between the test, and the use of
	 * the path.
	 */
}