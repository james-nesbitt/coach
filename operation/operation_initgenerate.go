package operation

import (
	"io"
	"os"
	"strings"

	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/initialize"
	"github.com/james-nesbitt/coach/log"
)

type InitGenerateOperation struct {
	log  log.Log
	conf *conf.Project

	force bool

	handler   string
	root      string   // Path to root
	skip      []string // skip some files
	sizeLimit int64

	output string // output file (or use logger to output to logger)
}

func (operation *InitGenerateOperation) Id() string {
	return "init"
}
func (operation *InitGenerateOperation) Flags(flags []string) bool {
	operation.root, _ = operation.conf.Paths.Path("project-root")
	operation.handler = "yaml"
	operation.output = "logger"
	operation.skip = []string{}
	operation.sizeLimit = 1024

	for index := 0; index < len(flags); index++ {
		flag := flags[index]
		switch flag {
		case "--test":
			operation.handler = "test"
			operation.output = "logger"
		case "-f":
			fallthrough
		case "--file":
			index++
			if index+1 > len(flags) {
				operation.output = "coachinit."
			} else if strings.HasPrefix(flags[index], "-") {
				operation.output = "coachinit."
			} else {
				operation.output = flags[index]
			}
		}
	}

	return true
}
func (operation *InitGenerateOperation) Help(topics []string) {
	operation.log.Message(`Operation: InitGenerate`)
}
func (operation *InitGenerateOperation) Run(logger log.Log) bool {
	logger.Info("running init operation:" + operation.output)

	var writer io.Writer
	switch operation.output {
	case "logger":
		fallthrough
	case "":
		writer = logger
	default:
		if strings.HasSuffix(operation.output, ".") {
			operation.output = operation.output + operation.handler
		}
		if fileWriter, err := os.Create(operation.output); err == nil {
			operation.skip = append(operation.skip, operation.output)
			writer = io.Writer(fileWriter)
			defer fileWriter.Close()
			logger.Message("Opening file for init generation output: " + operation.output)
		} else {
			logger.Error("Could not open output file to write init to:" + operation.output)
		}
	}

	initialize.Init_Generate(logger.MakeChild("init-generate"), operation.handler, operation.root, operation.skip, operation.sizeLimit, writer)

	return true
}
