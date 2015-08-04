package main

import (
	"path"

	"io/ioutil"
	"gopkg.in/yaml.v2"
)

const confsource_secrettokenyaml_filename string = "secrets.yml"

/**
 * Get more conf tokens from the secrets file, which is usually at .coach/secrets/secrets.yaml
 */
func (conf *Conf) from_SecretsYaml(log Log) {
	log.Debug(LOG_SEVERITY_DEBUG_LOTS,"Updating Conf from secrets")

	if confPath, ok := conf.Path("secrets"); ok {
		// get the path to where the config file should be
		confPath = path.Join(confPath, confsource_secrettokenyaml_filename)

		log.Debug(LOG_SEVERITY_DEBUG_WOAH,"Project secrets file:"+confPath)

		// read the config file
		yamlFile, err := ioutil.ReadFile(confPath)
		if err!=nil {
			log.Error("Could not read the YAML secrets file: "+err.Error())
			return
		}

		// replace tokens in the yamlFile
		yamlFile = []byte( conf.TokenReplace(string(yamlFile)) )
		log.Debug(LOG_SEVERITY_DEBUG_STAAAP,"YAML (tokenized):"+ string(yamlFile))

		// parse the config file contents as a ConfSource_secrettokensyaml object
		source := new(Conf_SecretsYaml)
		if err := yaml.Unmarshal([]byte(yamlFile), source); err!=nil {
			log.Error("YAML marshalling of the YAML secrets file failed: "+err.Error())
			return
		}
		log.DebugObject(LOG_SEVERITY_DEBUG_STAAAP,"secrets source:", *source)

		conf.Merge( source.toConf() )

	} else {
		log.Debug(LOG_SEVERITY_DEBUG_LOTS,"YAML secrets file not found, no project coach folder:"+confPath)
	}

}

/**
 * A ConfSource that comes from a Yaml source for secrets in the .coach/secrets path
 */
type Conf_SecretsYaml struct {
	Secrets map[string]string   `yaml:"Secrets,omitempty"` // secrettokens
}
func (source *Conf_SecretsYaml) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&source.Secrets)
}
func (source *Conf_SecretsYaml) toConf() Conf {
	conf := Conf{
		Tokens: map[string]string{},
	}

	if source.Secrets!=nil {
		for secret, value := range source.Secrets {
			conf.Tokens[secret] = value
		}
	}

	return conf
}

