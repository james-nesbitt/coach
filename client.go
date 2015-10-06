package main

import (
	"os"
	"strings"
	"path"
	"errors"

	"encoding/json" // we use this to convert Config and HostConfig to strings for token replacement.  It could easily be replaced.

	docker "github.com/fsouza/go-dockerclient"
)

var cache dockerClientCache

type dockerClientCache struct {
	client *docker.Client
	log Log
	images []docker.APIImages
	containers []docker.APIContainers
}

func GetClient(conf Conf, log Log) (*docker.Client, error) {

	var client *docker.Client
	var err error 

	log.DebugObject(LOG_SEVERITY_DEBUG_WOAH,"Docker client conf: ",conf.Docker)

	if (strings.HasPrefix(conf.Docker.Host, "tcp://")) {

		if _, err := os.Stat(conf.Docker.CertPath); err == nil {

			// TCP DOCKER CLIENT WITH CERTS
			client, err = docker.NewTLSClient(
				conf.Docker.Host,
				path.Join(conf.Docker.CertPath, "cert.pem"),
				path.Join(conf.Docker.CertPath, "key.pem"),
				path.Join(conf.Docker.CertPath, "ca.pem"),
			)

		} else {

			// TCP DOCKER CLIENT WITHOUT CERTS
			client, err = docker.NewClient(conf.Docker.Host)
		}

	} else if (strings.HasPrefix(conf.Docker.Host, "unix://")) {

		if _, err := os.Stat(conf.Docker.Host[7:]); err != nil {
			log.Fatal("Docker socket does not exist: ["+conf.Docker.Host+"] "+err.Error())
		} else {
			client, err = docker.NewClient(conf.Docker.Host)
		}

	} else {

		err = errors.New("Unknown client host :"+conf.Docker.Host)

	}

	if err != nil {
		log.Fatal(err.Error())
	}
	if client != nil {
		cache = dockerClientCache{ client:client, log:log }
	}

	return client, err
}

/**
 * Keep a global cache of retrieved images and containers
 * so that functionality that retrieves image and container
 * data doesn't need multiple trips.
 */


// reload image and container lists from the docker client
func (cache *dockerClientCache) refresh(refreshImages bool, refreshContainers bool) {

	var err error

	if refreshImages {
		filters := map[string][]string{}
		//filters["id"] = []string{"something"} // As it turns out, filters are actually kind of useless : https://stackoverflow.com/questions/24659300/how-to-use-docker-images-filter

		options := docker.ListImagesOptions{
			Filters: filters,
		}
		cache.images, err = cache.client.ListImages(options)

		if err != nil {
			cache.log.Fatal(err.Error())
		}
	}
	if refreshContainers {
		filters := map[string][]string{}
		//filters["id"] = []string{"project"} // As it turns out, filters are actually kind of useless : https://stackoverflow.com/questions/24659300/how-to-use-docker-images-filter
		// filters["status"] = []string{"running"}

		options := docker.ListContainersOptions{
			All: true,
			Filters: filters,
		}

		cache.containers, err = cache.client.ListContainers(options)

		if err != nil {
			cache.log.Fatal(err.Error())
		}
	}

}
// get a list of images from the docker image cache
func (cache *dockerClientCache) getImages(refresh bool) []docker.APIImages {
	if refresh || cache.images==nil {
		cache.refresh(true, false)
	}
	return cache.images
}
// get a list of containers from the docker container cache
func (cache *dockerClientCache) getContainers(refresh bool) []docker.APIContainers {
	if refresh || cache.containers==nil {
		cache.refresh(false, true)
	}
	return cache.containers	
}

/**
 * Get matching Image information
 */
func MatchImages(client *docker.Client, imagePrefix string) []docker.APIImages {
	images := []docker.APIImages{}

	for _, image := range cache.getImages(false) {
		eachimage:
			for _, tag := range image.RepoTags {
				if strings.HasPrefix(tag, imagePrefix) {
					images = append(images, image)
					continue eachimage
				}
			}
	}

	return images
}

func (node *Node) GetImages()  []docker.APIImages {
	return MatchImages(node.client, node.GetImageName())
}
func (node *Node) hasImage() bool {
	return len(node.GetImages())>0
}

/**
 * Get Node container information
 */

func MatchContainers(client *docker.Client, containerPrefix string, running bool) []docker.APIContainers {
	containers := []docker.APIContainers{}

	for _, container := range cache.getContainers(false) {
		if running && !strings.Contains(container.Status, "RUNNING") {
			continue
		}

		EachContainer:

		for _, containerName := range container.Names {
			if strings.Contains(containerName, "/"+containerPrefix) {
				containers = append(containers, container)
				continue EachContainer
			}
		}

	}

	return containers
}


func (node *Node) GetContainers(running bool) []docker.APIContainers {
	return MatchContainers(node.client, node.MachineName , running)
}
func (instance *Instance) GetContainer(running bool) (docker.APIContainers, bool) {
	matches := MatchContainers(instance.Node.client, instance.GetContainerName() , running)
	if len(matches)>0 {
		return matches[0], true
	} else {
		return docker.APIContainers{}, false
	}
}
func (instance *Instance) HasContainer(running bool) bool {
	_, found := instance.GetContainer(running)
	return found
}

/**
 * In cases where we transform HostConfigs, we need to start the transformation
 * from a copy, so as to leave the original untouched.  This occurs for each
 * instance config, which needs to have it's links and volumesFrom checked
 *
 * This code is exhaustive and silly, and should be refactored.  It is atomic
 * and can be changes withing the functions as needed, as long as a separate
 * copy of the object is returned with the same values.
 */

func DockerClient_Config_Copy(config docker.Config) docker.Config {
	newConfig := docker.Config{
		Hostname: config.Hostname,
		Domainname: config.Domainname,
		User: config.User,
		Memory: config.Memory,
		MemorySwap: config.MemorySwap,
		CPUShares: config.CPUShares,
		CPUSet: config.CPUSet,
		AttachStdin: config.AttachStdin,
		AttachStdout: config.AttachStdout,
		AttachStderr: config.AttachStderr,
		PortSpecs: config.PortSpecs,
		ExposedPorts: config.ExposedPorts,
		Tty: config.Tty,
		OpenStdin: config.OpenStdin,
		StdinOnce: config.StdinOnce,
		Env: config.Env,
		Cmd: config.Cmd,
		DNS: config.DNS,
		Image: config.Image,
		Volumes: config.Volumes,
		VolumesFrom: config.VolumesFrom,
		WorkingDir: config.WorkingDir,
		MacAddress: config.MacAddress,
		Entrypoint: config.Entrypoint,
		NetworkDisabled: config.NetworkDisabled,
		SecurityOpts: config.SecurityOpts,
		OnBuild: config.OnBuild,
		Labels: config.Labels,
	}
	return newConfig
}
/**
 * In cases where we transform HostConfigs, we need to start the transformation
 * from a copy, so as to leave the original untouched.
 *
 * @TODO refactor this to make it a little lesss ridiculous
 */
func DockerClient_HostConfig_Copy(config docker.HostConfig) docker.HostConfig {
	newConfig := docker.HostConfig{
// 		Binds: config.Binds,
		CapAdd: config.CapAdd,
		CapDrop: config.CapDrop,
		ContainerIDFile: config.ContainerIDFile,
		LxcConf: []docker.KeyValuePair{}, //config.LxcConf,
		Privileged: config.Privileged,
		PortBindings: map[docker.Port][]docker.PortBinding{}, //config.PortBindings,
// 		Links: config.Links,
		PublishAllPorts: config.PublishAllPorts,
		DNS: config.DNS,
		DNSSearch: config.DNSSearch,
// 		ExtraHosts: config.ExtraHosts,
// 		VolumesFrom: config.VolumesFrom,
		NetworkMode: config.NetworkMode,
		IpcMode: config.IpcMode,
		PidMode: config.PidMode,
		UTSMode: config.UTSMode,
		RestartPolicy: config.RestartPolicy,
		Devices: config.Devices,
		LogConfig: config.LogConfig,
		ReadonlyRootfs: config.ReadonlyRootfs,
		SecurityOpt: config.SecurityOpt,
		CgroupParent: config.CgroupParent,
		Memory: config.Memory,
		MemorySwap: config.MemorySwap,
		CPUShares: config.CPUShares,
		CPUSet: config.CPUSet,
		CPUQuota: config.CPUQuota,
		CPUPeriod: config.CPUPeriod,
		Ulimits: config.Ulimits,
	}

	for _, bind := range config.Binds {
		newConfig.Binds = append(newConfig.Binds, bind)
	}
	for _, conf := range config.LxcConf {
		newConfig.LxcConf = append(newConfig.LxcConf, conf)
	}
	for port, binding := range config.PortBindings {
		newConfig.PortBindings[port] = binding
	}
	for _, link := range config.Links {
		newConfig.Links = append(newConfig.Links, link)
	}
	for _, extraHost := range config.ExtraHosts {
		newConfig.ExtraHosts = append(newConfig.ExtraHosts, extraHost)
	}
	for _, volumesFrom := range config.VolumesFrom {
		newConfig.VolumesFrom = append(newConfig.VolumesFrom, volumesFrom)
	}

	return newConfig
}

/**
 * There are cases where we want to do token replacement on Config and
 * HostConfig objects, to get instance related tokens.  The easiest way
 * I could think of doing it was to Marshall the configs to json, do
 * a string replacement, and the UnMarshall them back to objects.
 *
 * So far it's pretty reliable, but it's probably not greate for performance.
 */

func ConfigTokenReplace(config docker.Config, tokens map[string]string) docker.Config {
	bytelist, _ := json.Marshal(config)
	jsonbytes := string(bytelist)

	for key, value := range tokens {
		jsonbytes = strings.Replace(jsonbytes, "%"+key, value, -1)
	}

	var newConfig docker.Config
	json.Unmarshal([]byte(jsonbytes), &newConfig)
  return newConfig
}
func HostConfigTokenReplace(config docker.HostConfig, tokens map[string]string) docker.HostConfig {
	bytelist, _ := json.Marshal(config)
	jsonbytes := string(bytelist)

	for key, value := range tokens {
		jsonbytes = strings.Replace(jsonbytes, "%"+key, value, -1)
	}

	var newConfig docker.HostConfig
	json.Unmarshal([]byte(jsonbytes), &newConfig)
	return newConfig
}
