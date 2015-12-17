package project

import (
	"os"
	"strings"

	"github.com/james-nesbitt/coach-tools/log"
)

/**
 * Get the docker client
 *
 * What we do is we test for the following env var, to know if a remote host is projectigured.
 * export DOCKER_HOST=tcp://192.168.59.103:2376
 *
 */
func (project *Project) from_DefaultDockerClient(clientLog log.Log) {
	if project.Docker.Host!="" {
		clientLog.Debug(LOG_SEVERITY_DEBUG_LOTS, "Project object already has docker client projectguration. Skipping default docker client setting.", nil)
		return
	}

	if DockerHost := os.Getenv("DOCKER_HOST"); DockerHost=="" {
		clientLog.Debug(LOG_SEVERITY_DEBUG, "No local environment DOCKER settings found, assuming a locally running docker client will be found.", nil)
		project.Docker.Host = "unix:///var/run/docker.sock"
		return
	} else {
		project.Docker.Host = DockerHost
	}

	// if we have no cert path, and we are going to use a TCP socket, test for a default cert path.
	if (project.Docker.CertPath=="" && strings.HasPrefix(project.Docker.Host, "tcp://")) {
		if DockerCertPath := os.Getenv("DOCKER_CERT_PATH"); DockerCertPath=="" {
			project.Docker.CertPath = DockerCertPath
		}
	}

	clientLog.Debug(LOG_SEVERITY_DEBUG_WOAH, "Default Docker client settings:", project.Docker)
}
