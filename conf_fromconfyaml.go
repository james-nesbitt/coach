package main

import (
	"path"

	"io/ioutil"
	"gopkg.in/yaml.v2"
)

const confsource_projectyaml_filename string = "conf.yml"

/**
 * Retrive conf settings from the coach file, which is usually .coach/conf.yml"
 */
func (conf *Conf) from_YamlConf(log Log) {
	log.Debug(LOG_SEVERITY_DEBUG_LOTS,"Updating Conf from YAML")

	if confPath, ok := conf.Path("projectcoach"); ok {
		// get the path to where the config file should be
		confPath = path.Join(confPath, confsource_projectyaml_filename)

		log.Debug(LOG_SEVERITY_DEBUG_WOAH,"Project coach file:"+confPath)

		// read the config file
		yamlFile, err := ioutil.ReadFile(confPath)
		if err!=nil {
			log.Error("Could not read the YAML file: "+err.Error())
			return
		}

		// replace tokens in the yamlFile
		yamlFile = []byte( conf.TokenReplace(string(yamlFile)) )
		log.Debug(LOG_SEVERITY_DEBUG_STAAAP,"YAML (tokenized):"+ string(yamlFile))

		// parse the config file contents as a ConfSource_projectyaml object
		source := new(Conf_Yaml)
		if err := yaml.Unmarshal(yamlFile, source); err!=nil {
			log.Error("YAML marshalling of the YAML conf file failed: "+err.Error())
			return
		}
		log.DebugObject(LOG_SEVERITY_DEBUG_STAAAP,"YAML source:", *source)

		conf.Merge( source.toConf(log) )

	} else {
		log.Debug(LOG_SEVERITY_DEBUG_LOTS,"YAML file not found, no project coach folder:"+confPath)
	}

}

/**
 * A ConfSource that comes from a Yaml source in the project .rodo folder
 */
type Conf_Yaml struct {
	Project string							`yaml:"Project,omitempty"`
	Author string								`yaml:"Author,omitempty"`

	Paths map[string]string			`yaml:"Paths,omitempty"`

	Tokens map[string]string		`yaml:"Tokens,omitempty"`

	Settings map[string]string	`yaml:"Settings,omitempty"`

	Docker	DockerClientConf		`yaml:"Docker,omitempty"`
}

func (source *Conf_Yaml) toConf(log Log) Conf {
	conf := Conf{
		Paths: map[string]string{},
		Tokens: map[string]string{},
	}

	if source.Project!="" {
		conf.Project = source.Project
	}
	if source.Author!="" {
// 		conf.Author = source.Author
	}

	if source.Paths!=nil {
		for index, path := range source.Paths {
			conf.Paths[index] = path
		}
	}
	if source.Tokens!=nil {
		for token, value := range source.Tokens {
			conf.Tokens[token] = value
		}
	}
	if source.Docker.Host!="" {
		conf.Docker.Host = source.Docker.Host
	}
	if source.Docker.CertPath!="" {
		conf.Docker.CertPath = source.Docker.CertPath
	}

log.DebugObject(LOG_SEVERITY_DEBUG_STAAAP,"YAML CONVERT:", source.Docker, conf.Docker)

	return conf
}
