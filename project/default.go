package project

func (project *Project) Default(projectLog Log) {
	projectLog.Debug(LOG_SEVERITY_DEBUG_LOTS,"Creating default Project")

	homeDir := "."
	if currentUser,  err := user.Current(); err==nil {
		homeDir = currentUser.HomeDir
	} else {
		homeDir = os.Getenv("HOME")
	}

	wd, _ := os.Getwd()
	_, err := os.Stat( path.Join(wd, coachProjectigFolder) )
	RootSearch:
		for err!=nil {
			wd = path.Dir(wd)
			if (wd==homeDir || wd=="." || wd=="/") {
				projectLog.Info("Could not find a project folder, coach will assume that this project is not initialized.")
				wd, _ = os.Getwd()
				break RootSearch
			}
			_, err = os.Stat(path.Join(wd, coachProjectigFolder) )
		}

	/**
	 * Set up some frequesntly used paths
	 */
	project.Paths["user-home"] = homeDir
	project.Paths["user-coach"] = path.Join(project.Paths["userhome"],coachProjectigFolder)
	project.Paths["project-root"] = wd
	project.Paths["project-coach"] = path.Join(wd,coachProjectigFolder)

}