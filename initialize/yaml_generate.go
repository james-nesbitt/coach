package initialize

import (
	"bufio"
	"io"
	"os"
	"path"
	"strings"

	"github.com/james-nesbitt/coach/log"
)

type YMLInitGenerator struct {
	output io.Writer
	logger log.Log
}

func (generator *YMLInitGenerator) generateSingleFile(fullPath string, sourcePath string) bool {
	singleFile, _ := os.Open(fullPath)
	defer singleFile.Close()

	generator.logger.Debug(log.VERBOSITY_DEBUG_LOTS, "GENERATE SINGLE FILE: ", singleFile.Name())
	generator.output.Write([]byte("- Type: File\n"))
	generator.output.Write([]byte("  Path: " + sourcePath + "\n"))
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
		generator.output.Write([]byte("    " + line))
	}
	return true
}
func (generator *YMLInitGenerator) generateGit(fullPath string, sourcePath string) bool {

	gitUrl := ""

	if configFile, err := os.Open(path.Join(fullPath, ".git", "config")); err == nil {

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
		generator.logger.Error("Could not open .git/config in " + sourcePath)
		return false
	}

	if gitUrl == "" {
		generator.logger.Error("Could not determine GIT Url from .git/config")
		return false
	}

	generator.logger.Debug(log.VERBOSITY_DEBUG_LOTS, "GENERATE GIT FILE: ", sourcePath)
	generator.output.Write([]byte("- Type: GitClone\n"))
	generator.output.Write([]byte("  Path: " + sourcePath + "\n"))
	generator.output.Write([]byte("  Url: " + gitUrl + "\n"))
	return true
}
