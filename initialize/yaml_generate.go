package initialize

import (
	"io"
	"bufio"
	"os"
	"path"
	"strings"

	"github.com/james-nesbitt/coach/log"
)

type YMLInitGenerator struct {
	output io.Writer
	logger log.Log
}
func (generator *YMLInitGenerator) Generate(logger log.Log, path string, skip []string, output io.Writer) bool {
	generator.logger = logger
	generator.output = output

	return generator.generateFileRecursive(path, "", skip)
}
func (generator *YMLInitGenerator)  generateFileRecursive(sourceRootPath string, sourcePath string, skip []string) bool {
	fullPath := sourceRootPath
	logger := generator.logger

	if sourcePath != "" {
		fullPath = path.Join(fullPath, sourcePath)
	}

	for _, skipEach := range skip {
		if skipEach == sourcePath {
			logger.Info("Skipping marked skip file :"+sourcePath)
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
			if !generator.generateFileRecursive(sourceRootPath, childSourcePath, skip) {
				logger.Warning("Resursive generate failed")
			}

		}

	} else {
		// add file copy
		if generator.generateSingleFile(fullPath, sourcePath) {
			logger.Info("Generated file (recursively): " + sourcePath)
			return true
		} else {
			logger.Warning("Failed to generate file: " + sourcePath)
			return false
		}
		return true
	}
	return true
}

func (generator *YMLInitGenerator) generateSingleFile(fullPath string, sourcePath string) bool {
	singleFile, _ := os.Open(fullPath)
	defer singleFile.Close()

	generator.logger.Debug(log.VERBOSITY_DEBUG_LOTS, "GENERATE SINGLE FILE: ", singleFile.Name())
	generator.output.Write([]byte("- Type: File\n"))
	generator.output.Write([]byte("  Path: "+sourcePath+"\n"))
	generator.output.Write([]byte("  Contents: |\n"))

	r := bufio.NewReader(singleFile)
	for {
		line, err := r.ReadString(10) // 0x0A separator = newline
		if err == io.EOF {
			break
		} else if err != nil {
			generator.logger.Error(err.Error())
			return false // if you return error
		}
		generator.output.Write([]byte("    "+line+"\n"))
	}
	return true
}
func (generator *YMLInitGenerator) generateGit(fullPath string, sourcePath string) bool {

	gitUrl := ""

	if configFile, err := os.Open(path.Join(fullPath, ".git", "config")); err==nil {

		r := bufio.NewReader(configFile)
		for {
			line, err := r.ReadString(10) // 0x0A separator = newline
			if err == io.EOF {
				break
			} else if err != nil {
				generator.logger.Error(err.Error())
				return false // if you return error
			}
			if strings.Contains(line, "url =") {
				lineSplit := strings.Split(line, "url =")
				gitUrl = strings.Trim(lineSplit[len(lineSplit)-1], " ")
				break
			}
		}

	} else {
		generator.logger.Error("Could not open .git/config in " +sourcePath)
		return false
	}

	if gitUrl=="" {
		generator.logger.Error("Could not determine GIT Url from .git/config")
		return false
	}

	generator.logger.Debug(log.VERBOSITY_DEBUG_LOTS, "GENERATE GIT FILE: ", sourcePath)
	generator.output.Write([]byte("- Type: GitClone\n"))
	generator.output.Write([]byte("  Path: "+sourcePath+"\n"))
	generator.output.Write([]byte("  Url: "+gitUrl+"\n"))
	return true	
}
