package initialize

import (
	"os"
	"os/exec"

	"github.com/james-nesbitt/coach/log"
)

func (tasks *InitTasks) Init_Git_Run(logger log.Log, source string) bool {

	if source == "" {
		logger.Error("You have not provided a git target $/> coach init git https://github.com/aleksijohansson/docker-drupal-coach")
		return false
	}

	url := source
	path := tasks.root

	cmd := exec.Command("git", "clone", "--progress", url, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = logger
	cmd.Stderr = logger

	err := cmd.Start()

	if err != nil {
		logger.Error("Failed to clone the remote repository [" + url + "] => " + err.Error())
		return false
	}

	logger.Message("Clone remote repository to local project folder [" + url + "]")
	err = cmd.Wait()

	if err != nil {
		logger.Error("Failed to clone the remote repository [" + url + "] => " + err.Error())
		return false
	}

	tasks.AddMessage("Cloned remote repository [" + url + "] to local project folder")
	tasks.AddFile(".coach/CREATEDFROM.md", `THIS PROJECT WAS CREATED FROM GIT`)

	return true
}

type InitTaskGitClone struct {
	InitTaskFileBase
	root string

	path string
	url  string
}

func (task *InitTaskGitClone) RunTask(logger log.Log) bool {
	if task.root == "" || task.url == "" {
		logger.Error("EMPTY ROOT PASSED TO GIT: " + task.root)
		return false
	}

	destinationPath := task.path
	url := task.url

	if !task.MakeDir(logger, destinationPath, false) {
		return false
	}

	destinationAbsPath, ok := task.absolutePath(destinationPath, true)
	if !ok {
		logger.Warning("Invalid copy destination path: " + destinationPath)
		return false
	}

	cmd := exec.Command("git", "clone", "--progress", url, destinationAbsPath)
	cmd.Stderr = logger
	err := cmd.Start()

	if err != nil {
		logger.Error("Failed to clone the remote repository [" + url + "] => " + err.Error())
		return false
	}

	err = cmd.Wait()

	if err != nil {
		logger.Error("Failed to clone the remote repository [" + url + "] => " + err.Error())
		return false
	}

	logger.Message("Cloned remote repository [" + url + "] to local path " + destinationPath)
	return true
}
