package main

import (
	"os"
	"strings"
	"path"
	"errors"

	"encoding/json" // we use this to convert Config and HostConfig to strings for token replacement.  It could easily be replaced.

	docker "github.com/fsouza/go-dockerclient"
)

func GetClient(conf Conf, log Log) (*docker.Client, error) {

	log.DebugObject(LOG_SEVERITY_DEBUG_WOAH,"Docker client conf: ",conf.Docker)

	if (strings.HasPrefix(conf.Docker.Host, "tcp://")) {

		if _, err := os.Stat(conf.Docker.CertPath); err == nil {

			// TCP DOCKER CLIENT WITH CERTS

			client, err := docker.NewTLSClient(
				conf.Docker.Host,
				path.Join(conf.Docker.CertPath, "cert.pem"),
				path.Join(conf.Docker.CertPath, "key.pem"),
				path.Join(conf.Docker.CertPath, "ca.pem"),
			)
			if err != nil {
				log.Fatal(err.Error())
				return nil, err
			}
			return client, nil

		} else {

			// TCP DOCKER CLIENT WITHOUT CERTS

			client, err := docker.NewClient(conf.Docker.Host)
			if err != nil {
				log.Fatal(err.Error())
				return nil, err
			}
			return client, nil

		}

	} else if (strings.HasPrefix(conf.Docker.Host, "unix://")) {

		if _, err := os.Stat(conf.Docker.Host[7:]); err != nil {
			log.Fatal("Docker socket does not exist: ["+conf.Docker.Host+"] "+err.Error())
			return nil, err
		}

		client, err := docker.NewClient(conf.Docker.Host)
		if err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		return client, nil

	} else {

		err := errors.New("Unknown client host :"+conf.Docker.Host)
		return nil, err

	}
}

/**
 * Get matching Image information
 */
func MatchImages(client *docker.Client, imagePrefix string) []docker.APIImages {
	images := []docker.APIImages{}

	filters := map[string][]string{}
	//filters["id"] = []string{"laird"} // As it turns out, filters are actually kind of useless : https://stackoverflow.com/questions/24659300/how-to-use-docker-images-filter

	options := docker.ListImagesOptions{
		Filters: filters,
	}

	if allImages, err := client.ListImages(options); err==nil {
		for _, image := range allImages {
			eachimage:
				for _, tag := range image.RepoTags {
					if strings.HasPrefix(tag, imagePrefix) {
						images = append(images, image)
						continue eachimage
					}
				}
		}
	} else {
		return nil
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

	filters := map[string][]string{}
	//filters["id"] = []string{"project"} // As it turns out, filters are actually kind of useless : https://stackoverflow.com/questions/24659300/how-to-use-docker-images-filter

	if running {
		filters["status"] = []string{"running"}
	}

	options := docker.ListContainersOptions{
		All: true,
		Filters: filters,
	}

	if allContainers, err := client.ListContainers(options); err==nil {

		for _, container := range allContainers {
			EachContainer:

			for _, containerName := range container.Names {
				if strings.Contains(containerName, "/"+containerPrefix) {

					if running {
						if strings.Contains(container.Status, "RUNNING") {
							containers = append(containers, container)
							continue EachContainer
						}
					} else {
						containers = append(containers, container)
						continue EachContainer
					}

				}
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
 * So far it's pretty reliable
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
