package project

import (
	"path"

	"io/ioutil"
	"gopkg.in/yaml.v2"
)

/**
 * Get more project tokens from the secrets file, which is usually at .coach/secrets/secrets.yaml
 */
func (project *Project) from_SecretsYaml(secretsLog log.Log, projectPathKey string) bool {
	secretsLog.Debug(LOG_SEVERITY_DEBUG_LOTS,"Updating Project from secrets", nil)

	if projectPath, ok := project.Path(projectPathKey); ok {
		// get the path to where the projectig file should be
		projectPath = path.Join(projectPath, "secrets.yml")

		secretsLog.Debug(LOG_SEVERITY_DEBUG_WOAH,"Secrets file:"+projectPath, nil)

		// read the projectig file
		yamlFile, err := ioutil.ReadFile(projectPath)
		if err!=nil {
			secretsLog.Info("Could not read the YAML file ["+projectPath+"]: "+err.Error())
			return false
		}

		// replace tokens in the yamlFile
		yamlFile = []byte( project.TokenReplace(string(yamlFile)) )
		secretsLog.Debug(LOG_SEVERITY_DEBUG_STAAAP,"YAML (tokenized):"+ string(yamlFile))

		// parse the projectig file contents as a ProjectSource_secrettokensyaml object
		source := new(Project_SecretsYaml)
		if err := yaml.Unmarshal([]byte(yamlFile), source); err!=nil {
			secretsLog.Warning("YAML marshalling of the YAML project file failed ["+projectPath+"]: "+err.Error())
			return false
		}
		secretsLog.Debug(LOG_SEVERITY_DEBUG_STAAAP,"secrets source:", *source)

		project.Merge( source.toProject() )

		return true

	} else {
		secretsLog.Info(LOG_SEVERITY_DEBUG_LOTS,"YAML secrets file not found, no project coach folder:"+projectPath)
	}

	return false
}

/**
 * A ProjectSource that comes from a Yaml source for secrets in the .coach/secrets path
 */
type Project_SecretsYaml struct {
	Secrets map[string]string   `yaml:"Secrets,omitempty"` // secrettokens
}
func (source *Project_SecretsYaml) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&source.Secrets)
}
func (source *Project_SecretsYaml) toProject() Project {
	project := Project{
		Tokens: map[string]string{},
	}

	if source.Secrets!=nil {
		for secret, value := range source.Secrets {
			project.Tokens[secret] = value
		}
	}

	return project
}

