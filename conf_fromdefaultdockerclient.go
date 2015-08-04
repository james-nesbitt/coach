package main

import (
	"os"
	"strings"
)


/**
 * Get the docker client
 *
 * What we do is we test for the following env var, to know if a remote host is configured.
 * export DOCKER_HOST=tcp://192.168.59.103:2376
 *
 */
func (conf *Conf) from_DefaultDockerClient(log Log) {
	if conf.Docker.Host!="" {
		log.Debug(LOG_SEVERITY_DEBUG_LOTS, "Conf object already has docker client confguration. Skipping default docker client setting.")
		return
	}

	if DockerHost := os.Getenv("DOCKER_HOST"); DockerHost=="" {
		log.Debug(LOG_SEVERITY_DEBUG, "No local environment DOCKER settings found, assuming a locally running docker client will be found.")
		conf.Docker.Host = "unix:///var/run/docker.sock"
		return
	} else {
		conf.Docker.Host = DockerHost
	}

	// if we have no cert path, and we are going to use a TCP socket, test for a default cert path.
	if (conf.Docker.CertPath=="" && strings.HasPrefix(conf.Docker.Host, "tcp://")) {
		if DockerCertPath := os.Getenv("DOCKER_CERT_PATH"); DockerCertPath=="" {
			conf.Docker.CertPath = DockerCertPath
		}
	}

	log.DebugObject(LOG_SEVERITY_DEBUG_WOAH, "Default Docker client settings:", conf.Docker)
}
