package initialize

import (
	"io"

	"github.com/james-nesbitt/coach/log"
)

func Init_Generate(logger log.Log, handler string, path string, skip []string, output io.Writer) bool {
	logger.Message("GENERATING INIT")

	var generator Generator
	switch handler {
	case "yaml":
		generator = Generator(&YMLInitGenerator{})
	default:
		logger.Error("Unknown init generator (handler) "+handler)
		return false;
	}

	success := generator.Generate(logger, path, skip, output)
	if (success) {
		logger.Message("FINISHED GENERATING YML INIT")
	} else {
		logger.Error("ERROR OCCURRED GENERATING YML INIT")
	}
	return success
}

type Generator interface {

	Generate(logger log.Log, path string, skip []string, output io.Writer) bool

}
