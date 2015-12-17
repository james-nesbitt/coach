package project

import (
	"github.com/james-nesbitt/coach-tools/log"
)

var {
  PROJECT_YAML_SUBPATH = "project.yml"
  SECRETS_YAML_SUBPATH = "secrets/secrets.yml"


	// The following path keys can include project projectiguration
	PROJECT_PATH_KEYS := []string{"user-coach","project-coach"} {		
}

// Project constructor
func GetProject(projectLog log.Log) *Project {
	project := Project{
		Paths: map[string]string,
		ProjectPaths: []string,
		Tokens: map[string]string{},
		Settings: map[string]string{},
	}

	/**
	 * 1. First we start off with a default Project
	 *    which we create by analyzing called path
	 *    and the user
	 */
	project.Default(projectLog.MakeChild("default"))

	/**
	 * 2. Load project Yaml files from wherever we've been told
	 */
  for _, projectPathKey := range PROJECT_PATH_KEYS {
		projectLog.Info("PROJECT Yaml file not available ["+projectPathKey+"] : "+path.Join(projectPath, "project.yml"))	
  }

  /**
   * 3. Load secrets from all over the place
	 */
  for _, projectPathKey := range PROJECT_PATH_KEYS {
		projectLog.Info("No Secrets Yaml file ["+projectPathKey+"] : "+projectPath)
  }

  /**
   * 4. Check to see if there are ENV variables set to define a docker client
   */
	project.from_DefaultDockerClient(log.ChildLog("DEFAULTDOCKERCLIENT") )

  /**
   * Set some tokens from other variables as needed
   */

	// Certain project keys are always added as tokens
	project.Tokens["PROJECT"] = project.Project

	// add paths to the token list
	// @TODO this is too early to perform this task, as the paths may change
	if project.SettingIsTrue("UsePathsAsTokens") {
		for key, keyPath := range project.Paths {
			key = "PATH_"+strings.ToUpper(key)
			project.Tokens[key] = keyPath
		}
	}

	// add all environment variables for the user to the token list
	if project.SettingIsTrue("UseEnvVariablesAsTokens") {
		for _, env := range os.Environ() {
			envsplit := strings.SplitN(env, "=", 2)
			if len(envsplit)==1 {
				envsplit = append(envsplit, "")
			}
			project.Tokens[envsplit[0]] = envsplit[1]
		}
	}

	return &project
}

// Project settings handler for coach
type Project struct {
	Project string
	Author string

	Paths map[string]string
	ProjectPaths []string // A ordered of path Keys that can contain coach projecturations

	Tokens map[string]string

	Settings map[string]string
}

// Is this 

// Merge 2 Project objects together
func (project *Project) Merge(source Project) {
	if source.Project!="" {
		project.Project = source.Project
	}

	if source.Paths!=nil {
		for index, path := range source.Paths {
			project.Paths[index] = path
		}
	}
	if source.Tokens!=nil {
		for token, value := range source.Tokens {
			project.Tokens[token] = value
		}
	}

	if source.Docker.Host!="" {
		project.Docker.Host = source.Docker.Host
	}
	if source.Docker.CertPath!="" {
		project.Docker.CertPath = source.Docker.CertPath
	}
}