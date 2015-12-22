package conf

const (
	COACH_PROJECT_FLAG_DEFAULT_UsePathsAsTokens = true
	COACH_PROJECT_FLAG_DEFAULT_UseEnvVariablesAsTokens = false
)

func MakeProjectFlags() ProjectFlags {
	return ProjectFlags{
			UsePathsAsTokens: COACH_PROJECT_FLAG_DEFAULT_UsePathsAsTokens,
			UseEnvVariablesAsTokens: COACH_PROJECT_FLAG_DEFAULT_UseEnvVariablesAsTokens,
	}
}

// ProjectFlags are Configuration flags for the project settings
type ProjectFlags struct {
	UsePathsAsTokens bool// Should all of the paths be available as tokens
	UseEnvVariablesAsTokens bool // Should the running user ENV variables be available as tokens
}
