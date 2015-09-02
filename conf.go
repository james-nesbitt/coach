package main

import (
	"os"
	"os/user"

	"strings"
	"path"
)

const (
	// this is the folder that marks a coach project
	coachConfigFolder = ".coach"

	// Tokens get wrapped in these characters
	conf_token_wrapper_prefix string = "%"
)

// get a new conf object, configured from various sources
func GetConf(log Log) Conf {
	conf := Conf{
		Paths: map[string]string{},
		Tokens: map[string]string{},
		Targets: []string{},
	}

	conf.from_Default(false, log.ChildLog("DEFAULT"))

	conf.from_YamlConf(log.ChildLog("USER"), "usercoach")
	if !conf.from_YamlConf(log.ChildLog("PROJECT"), "projectcoach") {
		log.Warning("This project contains no CONF Yaml file")
	}

	log.DebugObject(LOG_SEVERITY_DEBUG_LOTS,"Docker client conf: ",conf)

	conf.from_SecretsYaml(log.ChildLog("SECRETS"), "usersecrets")
	conf.from_SecretsYaml(log.ChildLog("SECRETS"), "projectsecrets")

	conf.from_DefaultDockerClient(log.ChildLog("DEFAULTDOCKERCLIENT") )

	conf.Tokens["PROJECT"] = conf.Project

	return conf
}

type DockerClientConf struct {
	Host			string			`json:"Host,omitempty" yaml:"Host,omitempty"`
	CertPath	string			`json:"CertPath,omitempty" yaml:"CertPath,omitempty"`
}

// Conf object
type Conf struct {
	Project string
	Author string

	Paths map[string]string

	Tokens map[string]string

	Targets []string

	Docker	DockerClientConf
}

func (conf *Conf) Merge(source Conf) {
	if source.Project!="" {
		conf.Project = source.Project
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
}

// Replace all occurences of tokens in a string with their token values
// @TODO for a large number of keys, it will be faster to not iterate through all the keys
func TokenMapReplace(tokens map[string]string, value string) string {
	for key, token := range tokens {
		value = strings.Replace(value, conf_token_wrapper_prefix+key, token, -1)
	}
	return value
}
func (conf *Conf) TokenReplace(value string) (string) {
	return TokenMapReplace(conf.Tokens, value)
}
// return a path for a key it if is set
func (conf *Conf) Path(name string) (string, bool) {
	if path, ok := conf.Paths[name]; ok {
		return path, ok
	} else {
		return "", ok
	}
}

func (conf *Conf) from_Default(includeEnv bool, log Log) {
	log.Debug(LOG_SEVERITY_DEBUG_LOTS,"Creating default Conf")

	homeDir := "."
	if currentUser,  err := user.Current(); err==nil {
		homeDir = currentUser.HomeDir
	} else {
		homeDir = os.Getenv("HOME")
	}

	wd, _ := os.Getwd()
	_, err := os.Stat( path.Join(wd, coachConfigFolder) )
	RootSearch:
		for err!=nil {
			wd = path.Dir(wd)
			if (wd==homeDir || wd=="." || wd=="/") {
				log.Warning("Could not find a project folder, coach will assume that this project is not initialized.")
				break RootSearch
			}
			_, err = os.Stat(path.Join(wd, coachConfigFolder) )
		}

	/**
	 * 1. First we start off with a default Conf
	 */
	conf.Paths["userhome"] = homeDir
	conf.Paths["usercoach"] = path.Join(conf.Paths["userhome"],coachConfigFolder)
	conf.Paths["usertemplates"] = path.Join(conf.Paths["usercoach"],"templates")
	conf.Paths["usersecrets"] = path.Join(conf.Paths["usercoach"],"secrets")
	conf.Paths["project"] = wd
	conf.Paths["projectcoach"] = path.Join(wd,coachConfigFolder)
	conf.Paths["projectsecrets"] = path.Join(conf.Paths["projectcoach"],"secrets") // keep secret things in one place for gitignore
	conf.Paths["build"] = conf.Paths["projectcoach"] // maybe for remote builds, this should be different?

	// add all environment variables for the user to the env list
	if includeEnv {
		for _, env := range os.Environ() {
			envsplit := strings.SplitN(env, "=", 2)
			if len(envsplit)==1 {
				envsplit = append(envsplit, "")
			}
			conf.Tokens[envsplit[0]] = envsplit[1]
		}
	}

}
