package initialize

import (
	"io"
	"os"
	"path"

	"github.com/james-nesbitt/coach/log"
)

func Init_Generate(logger log.Log, handler string, path string, skip []string, sizeLimit int64, output io.Writer) bool {
	logger.Message("GENERATING INIT")

	var generator Generator
	switch handler {
	case "yaml":
		generator = Generator(&YMLInitGenerator{logger: logger, output: output})
	default:
		logger.Error("Unknown init generator (handler) " + handler)
		return false
	}

	iterator := GenerateIterator{
		logger:    logger,
		output:    output,
		skip:      skip,
		sizeLimit: sizeLimit,
		generator: generator,
	}

	if iterator.Generate(path) {
		logger.Message("FINISHED GENERATING YML INIT")
		return true
	} else {
		logger.Error("ERROR OCCURRED GENERATING YML INIT")
		return false
	}
}

type GenerateIterator struct {
	logger log.Log
	output io.Writer

	skip      []string
	sizeLimit int64

	generator Generator
}

func (iterator *GenerateIterator) Generate(path string) bool {
	return iterator.generate_Recursive(path, "")
}
func (iterator *GenerateIterator) generate_Recursive(sourceRootPath string, sourcePath string) bool {
	fullPath := sourceRootPath

	if sourcePath != "" {
		fullPath = path.Join(fullPath, sourcePath)
	}

	for _, skipEach := range iterator.skip {
		if skipEach == sourcePath {
			iterator.logger.Info("Skipping marked skip file :" + sourcePath)
			return true
		}
	}

	// get properties of source dir
	info, err := os.Stat(fullPath)
	if err != nil {
		// @TODO do something log : source doesn't exist
		iterator.logger.Warning("File does not exist :" + fullPath)
		return false
	}

	mode := info.Mode()
	if mode.IsDir() {

		// check for GIT folder
		if _, err := os.Open(path.Join(fullPath, ".git")); err == nil {
			if iterator.generator.generateGit(fullPath, sourcePath) {
				iterator.logger.Info("Generated git file: " + sourcePath)
				return true
			} else {
				iterator.logger.Warning("Failed to generate git file: " + sourcePath)
			}
		}

		directory, _ := os.Open(fullPath)
		defer directory.Close()
		objects, err := directory.Readdir(-1)

		if err != nil {
			// @TODO do something log : source doesn't exist
			iterator.logger.Warning("Could not open directory")
			return false
		}

		for _, obj := range objects {

			//childSourcePath := source + "/" + obj.Name()
			childSourcePath := path.Join(sourcePath, obj.Name())
			if !iterator.generate_Recursive(sourceRootPath, childSourcePath) {
				iterator.logger.Warning("Resursive generate failed")
			}

		}

	} else if mode.IsRegular() {
		// add file copy
		if iterator.generator.generateSingleFile(fullPath, sourcePath) {
			iterator.logger.Info("Generated file (recursively): " + sourcePath)
			return true
		} else {
			iterator.logger.Warning("Failed to generate file: " + sourcePath)
			return false
		}
		return true
	} else {
		iterator.logger.Warning("Skipped generation non-regular file: " + sourcePath)
	}

	return true
}

type Generator interface {
	generateSingleFile(fullPath string, sourcePath string) bool
	generateGit(fullPath string, sourcePath string) bool
}
