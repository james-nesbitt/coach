package initialize

import (
	"io"
	"os"
	"path"

	"github.com/james-nesbitt/coach/log"
)

func Init_Generate(logger log.Log, handler string, path string, skip []string, output io.Writer) bool {
	logger.Message("GENERATING INIT")

	var generator Generator
	switch handler {
	case "yaml":
		generator = Generator(&YMLInitGenerator{logger: logger, output: output})
	default:
		logger.Error("Unknown init generator (handler) " + handler)
		return false
	}

	success := Init_Generate_Recursive(logger, generator, path, "", skip)
	if success {
		logger.Message("FINISHED GENERATING YML INIT")
	} else {
		logger.Error("ERROR OCCURRED GENERATING YML INIT")
	}
	return success
}

func Init_Generate_Recursive(logger log.Log, generator Generator, sourceRootPath string, sourcePath string, skip []string) bool {
	fullPath := sourceRootPath

	if sourcePath != "" {
		fullPath = path.Join(fullPath, sourcePath)
	}

	for _, skipEach := range skip {
		if skipEach == sourcePath {
			logger.Info("Skipping marked skip file :" + sourcePath)
			return true
		}
	}

	// get properties of source dir
	info, err := os.Stat(fullPath)
	if err != nil {
		// @TODO do something log : source doesn't exist
		logger.Warning("File does not exist :" + fullPath)
		return false
	}

	mode := info.Mode()
	if mode.IsDir() {

		// check for GIT folder
		if _, err := os.Open(path.Join(fullPath, ".git")); err == nil {
			if generator.generateGit(fullPath, sourcePath) {
				logger.Info("Generated git file: " + sourcePath)
				return true
			} else {
				logger.Warning("Failed to generate git file: " + sourcePath)
			}
		}

		directory, _ := os.Open(fullPath)
		defer directory.Close()
		objects, err := directory.Readdir(-1)

		if err != nil {
			// @TODO do something log : source doesn't exist
			logger.Warning("Could not open directory")
			return false
		}

		for _, obj := range objects {

			//childSourcePath := source + "/" + obj.Name()
			childSourcePath := path.Join(sourcePath, obj.Name())
			if !Init_Generate_Recursive(logger, generator, sourceRootPath, childSourcePath, skip) {
				logger.Warning("Resursive generate failed")
			}

		}

	} else if mode.IsRegular() {
		// add file copy
		if generator.generateSingleFile(fullPath, sourcePath) {
			logger.Info("Generated file (recursively): " + sourcePath)
			return true
		} else {
			logger.Warning("Failed to generate file: " + sourcePath)
			return false
		}
		return true
	} else {
		logger.Warning("Skipped generation non-regular file: " + sourcePath)
	}

	return true
}

type Generator interface {
	generateSingleFile(fullPath string, sourcePath string) bool
	generateGit(fullPath string, sourcePath string) bool
}
