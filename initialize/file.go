package initialize

import (
	"io"
	"os"
	"os/user"
	"path"
	"strings"

	"io/ioutil"
	"net/http"

	"github.com/james-nesbitt/coach/log"
)

type InitTaskFileBase struct {
	root string
}

func (task *InitTaskFileBase) absolutePath(targetPath string, addRoot bool) (string, bool) {
	if strings.HasPrefix(targetPath, "~") {
		return path.Join(task.userHomePath(), targetPath[1:]), !addRoot

		// I am not sure how reliable this function is
		// you passed an absolute path, so I can't add the root
		// } else if path.isAbs(targetPath) {
		// 	return targetPath, !addRoot

		// you passed a relative path, and want me to add the root
	} else if addRoot {
		return path.Join(task.root, targetPath), true

		// you passed path and don't want the root added (but is it already abs?)
	} else if targetPath != "" {
		return targetPath, true

		// you passed an empty string, and don't want the root added?
	} else {
		return targetPath, false
	}
}
func (task *InitTaskFileBase) userHomePath() string {
	if currentUser, err := user.Current(); err == nil {
		return currentUser.HomeDir
	} else {
		return os.Getenv("HOME")
	}
}

func (task *InitTaskFileBase) MakeDir(logger log.Log, makePath string, pathIsFile bool) bool {
	if makePath == "" {
		return true // it's already made
	}

	if pathDirectory, ok := task.absolutePath(makePath, true); !ok {
		logger.Warning("Invalid directory path: " + pathDirectory)
		return false
	}
	pathDirectory := path.Join(task.root, makePath)
	if pathIsFile {
		pathDirectory = path.Dir(pathDirectory)
	}

	if err := os.MkdirAll(pathDirectory, 0777); err != nil {
		// @todo something log
		return false
	}
	return true
}
func (task *InitTaskFileBase) MakeFile(logger log.Log, destinationPath string, contents string) bool {
	if !task.MakeDir(logger, destinationPath, true) {
		// @todo something log
		return false
	}

	if destinationPath, ok := task.absolutePath(destinationPath, true); !ok {
		logger.Warning("Invalid file destination path: " + destinationPath)
		return false
	}

	fileObject, err := os.Create(destinationPath)
	defer fileObject.Close()
	if err != nil {
		// @todo something log
		return false
	}
	if _, err := fileObject.WriteString(contents); err != nil {
		// @todo something log
		return false
	}

	return true
}

func (task *InitTaskFileBase) CopyFile(logger log.Log, destinationPath string, sourcePath string) bool {
	if destinationPath == "" || sourcePath == "" {
		logger.Warning("empty source or destination passed for copy")
		return false
	}

	sourcePath, ok := task.absolutePath(sourcePath, false)
	if !ok {
		logger.Warning("Invalid copy source path: " + sourcePath)
		return false
	}
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		logger.Warning("could not copy file as it does not exist [" + sourcePath + "] : " + err.Error())
		return false
	}
	defer sourceFile.Close()

	if !task.MakeDir(logger, destinationPath, true) {
		// @todo something log
		logger.Warning("could not copy file as the path to the destination file could not be created [" + destinationPath + "]")
		return false
	}
	destinationAbsPath, ok := task.absolutePath(destinationPath, true)
	if !ok {
		logger.Warning("Invalid copy destination path: " + destinationPath)
		return false
	}

	destinationFile, err := os.Open(destinationAbsPath)
	if err == nil {
		logger.Warning("could not copy file as it already exists [" + destinationPath + "]")
		destinationFile.Close()
		return false
	}

	destinationFile, err = os.Create(destinationAbsPath)
	if err != nil {
		logger.Warning("could not copy file as destination file could not be created [" + destinationPath + "] : " + err.Error())
		return false
	}

	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)

	if err == nil {
		sourceInfo, err := os.Stat(sourcePath)
		if err == nil {
			err = os.Chmod(destinationPath, sourceInfo.Mode())
			return true
		} else {
			logger.Warning("could not copy file as destination file could not be created [" + destinationPath + "] : " + err.Error())
			return false
		}
	} else {
		logger.Warning("could not copy file as copy failed [" + destinationPath + "] : " + err.Error())
	}

	return true
}

func (task *InitTaskFileBase) CopyRemoteFile(logger log.Log, destinationPath string, sourcePath string) bool {
	if destinationPath == "" || sourcePath == "" {
		return false
	}

	response, err := http.Get(sourcePath)
	if err != nil {
		logger.Warning("Could not open remote URL: " + sourcePath)
		return false
	}
	defer response.Body.Close()

	sourceContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Warning("Could not read remote file: " + sourcePath)
		return false
	}

	return task.MakeFile(logger, destinationPath, string(sourceContent))
}

func (task *InitTaskFileBase) CopyFileRecursive(logger log.Log, path string, source string) bool {
	sourceAbsPath, ok := task.absolutePath(source, false)
	if !ok {
		logger.Warning("Couldn't find copy source " + source)
		return false
	}
	return task.copyFileRecursive(logger, path, sourceAbsPath, "")
}
func (task *InitTaskFileBase) copyFileRecursive(logger log.Log, destinationRootPath string, sourceRootPath string, sourcePath string) bool {
	fullPath := sourceRootPath

	if sourcePath != "" {
		fullPath = path.Join(fullPath, sourcePath)
	}

	// get properties of source dir
	info,
		err := os.Stat(fullPath)
	if err != nil {
		// @TODO do something log : source doesn't exist
		logger.Warning("File does not exist :" + fullPath)
		return false
	}

	mode := info.Mode()
	if mode.IsDir() {

		directory, _ := os.Open(fullPath)
		objects, err := directory.Readdir(-1)

		if err != nil {
			// @TODO do something log : source doesn't exist
			logger.Warning("Could not open directory")
			return false
		}

		for _, obj := range objects {

			//childSourcePath := source + "/" + obj.Name()
			childSourcePath := path.Join(sourcePath, obj.Name())
			if !task.copyFileRecursive(logger, destinationRootPath, sourceRootPath, childSourcePath) {
				logger.Warning("Resursive copy failed")
			}

		}

	} else {
		// add file copy
		destinationPath := path.Join(destinationRootPath, sourcePath)
		if task.CopyFile(logger, destinationPath, sourceRootPath) {
			logger.Info("--> Copied file (recursively): " + sourcePath + " [from " + sourceRootPath + "]")
			return true
		} else {
			logger.Warning("--> Failed to copy file: " + sourcePath + " [from " + sourceRootPath + "]")
			return false
		}
		return true
	}
	return true
}
// perform a string replace on file contents
func (task *InitTaskFileBase) FileStringReplace(logger log.Log, targetPath string, oldString string, newString string, replaceCount int) bool {

	targetPath, ok := task.absolutePath(targetPath, false)
	if !ok {
		logger.Warning("Invalid string replace path: " + targetPath)
		return false
	}

	contents, err := ioutil.ReadFile(targetPath)
	if err != nil {
        logger.Error(err.Error())
	}

	contents = []byte( strings.Replace(string(contents), oldString, newString, replaceCount) )

	err = ioutil.WriteFile(targetPath, contents, 0644)
	if err != nil {
        logger.Error(err.Error())
	}
	return true
}

type InitTaskFile struct {
	InitTaskFileBase
	root string

	path     string
	contents string
}

func (task *InitTaskFile) RunTask(logger log.Log) bool {
	if task.path == "" {
		return false
	}

	if task.MakeFile(logger, task.path, task.contents) {
		logger.Message("--> Created file : " + task.path)
		return true
	} else {
		logger.Warning("--> Failed to create file : " + task.path)
		return false
	}
}

type InitTaskRemoteFile struct {
	InitTaskFileBase
	root string

	path string
	url  string
}

func (task *InitTaskRemoteFile) RunTask(logger log.Log) bool {
	if task.path == "" || task.root == "" || task.url == "" {
		return false
	}

	if task.CopyRemoteFile(logger, task.path, task.url) {
		logger.Message("--> Copied remote file : " + task.url + " -> " + task.path)
		return true
	} else {
		logger.Warning("--> Failed to copy remote file : " + task.url)
		return false
	}
}

type InitTaskFileCopy struct {
	InitTaskFileBase
	root string

	path   string
	source string
}

func (task *InitTaskFileCopy) RunTask(logger log.Log) bool {
	if task.path == "" || task.root == "" || task.source == "" {
		return false
	}

	if task.CopyFileRecursive(logger, task.path, task.source) {
		logger.Message("--> Copied file : " + task.source + " -> " + task.path)
		return true
	} else {
		logger.Warning("--> Failed to copy file : " + task.source + " -> " + task.path)
		return false
	}
}

type InitTaskFileStringReplace struct {
	InitTaskFileBase
	root string

	path string
	oldString  string
	newString	 string
	replaceCount 	 int
}

func (task *InitTaskFileStringReplace) RunTask(logger log.Log) bool {
	if task.path == "" || task.root == "" || task.oldString == "" || task.newString == "" {
		return false
	}
	if task.replaceCount==0 {
		task.replaceCount = -1
	}

	if task.FileStringReplace(logger, task.path, task.oldString, task.newString, task.replaceCount) {
		logger.Message("--> performed string replace on file : " + task.path)
		return true
	} else {
		logger.Warning("--> Failed to perform string replace on file : " + task.path)
		return false
	}
}
