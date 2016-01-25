package conf

import (
	"os"
	"strings"

	"github.com/james-nesbitt/coach/log"
)

const (
	COACH_PROJECT_FLAG_DEFAULT_UsePathsAsTokens        = true
	COACH_PROJECT_FLAG_DEFAULT_UseEnvVariablesAsTokens = false
)

func MakeProjectFlags() Flags {
	return Flags{
		UsePathsAsTokens:        COACH_PROJECT_FLAG_DEFAULT_UsePathsAsTokens,
		UseEnvVariablesAsTokens: COACH_PROJECT_FLAG_DEFAULT_UseEnvVariablesAsTokens,
	}
}

// ProjectFlags are Configuration flags for the project settings
type Flags struct {
	UsePathsAsTokens        bool // Should all of the paths be available as tokens
	UseEnvVariablesAsTokens bool // Should the running user ENV variables be available as tokens
}

func (project *Project) ProcessFlags(logger log.Log) {

	// add any paths as tokens if told to do so
	if project.Flags.UsePathsAsTokens {
		for _, pathKey := range project.Paths.PathOrder() {
			path, _ := project.Paths.Path(pathKey)
			project.Tokens.SetToken("PATH_"+strings.ToUpper(pathKey), path)
		}
	}

	// add all environment variables for the user to the token list
	if project.Flags.UseEnvVariablesAsTokens {
		for _, env := range os.Environ() {
			envsplit := strings.SplitN(env, "=", 2)
			if len(envsplit) == 1 {
				envsplit = append(envsplit, "")
			}
			project.Tokens[envsplit[0]] = envsplit[1]
		}
	}

}
