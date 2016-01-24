package libs

/**
 * @file Coach Client and ClientFactory, based on using an instances
 * of the FSouza docker remoteAPI client: https://github.com/fsouza/go-dockerclient
 *
 * COMPONENTS:
 *
 * FSouzaWrapper : a wrapper for the FSouza client, which caches image and container
 *   lists, as the remoteAPI does a poor job of filtering, and the retrieval process can
 *   take a long time over a remote connection.  This makes tasks such as confirming
 *   that a node has an image, or an instance has a container much faster.
 *
 * ClientFactory & ClientFactorySettings : The Coach Client factory, and settings
 * 		that create Client objects from
 */

import (
	"errors"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"

	"encoding/json"

	docker "github.com/fsouza/go-dockerclient"

	"github.com/james-nesbitt/coach/conf"
	"github.com/james-nesbitt/coach/log"
)

var (
	actionCache map[string]bool
)

func init() {
	actionCache = map[string]bool{}
}

/**
 * Coach: ClientFactory
 */
type FSouza_ClientFactorySettings struct {
	Host     string `json:"Host,omitempty" yaml:"Host,omitempty"`
	CertPath string `json:"CertPath,omitempty" yaml:"CertPath,omitempty"`
}

func (settings *FSouza_ClientFactorySettings) Settings() interface{} {
	return settings
}

type FSouza_ClientFactory struct {
	settings FSouza_ClientFactorySettings
	log      log.Log
	conf     *conf.Project
	client   *FSouza_Wrapper
}

// Provide a unique string identifier for this client-factory/type
func (clientFactory *FSouza_ClientFactory) Id() string {
	return "Docker:FSouza"
}

// Return a boolean for if this client factory matches the requirements
func (clientFactory *FSouza_ClientFactory) Match(requirements FactoryMatchRequirements) bool {
	clientFactory.log.Debug(log.VERBOSITY_DEBUG_STAAAP, "Match test for FSouza client factory:", requirements.Type, (requirements.Type == "docker"))
	return requirements.Type == "docker" || requirements.ID == clientFactory.Id() || requirements.Class == "FSouza_ClientFactory"
}

func (clientFactory *FSouza_ClientFactory) Init(logger log.Log, project *conf.Project, settings ClientFactorySettings) bool {
	clientFactory.log = logger
	clientFactory.conf = project

	// make sure that the settings that were given, where the proper "FSouza_ClientFactory" type
	typedSettings := settings.Settings()
	switch asserted := typedSettings.(type) {
	case *FSouza_ClientFactorySettings:
		clientFactory.settings = *asserted
	default:
		logger.Error("Invalid settings type passed to Fsouza Factory")
		logger.Debug(log.VERBOSITY_DEBUG, "Settings passed:", asserted)
	}

	// if we haven't made an actual fsouza docker client, then do it now
	if clientFactory.client == nil {
		if client, pk := clientFactory.makeFsouzaClientWrapper(logger.MakeChild("fsouza")); pk {
			clientFactory.client = client
			return true
		} else {
			logger.Error("Failed to create actual FSouza Docker client from client factory configuration")
			return false
		}
	}
	return true
}

// Actually build a Client-Wrapper object, one per factory, which makes sense because that is where the settings are
func (clientFactory *FSouza_ClientFactory) makeFsouzaClientWrapper(logger log.Log) (*FSouza_Wrapper, bool) {
	wrapper := &FSouza_Wrapper{}
	return wrapper, wrapper.Init(logger, clientFactory.settings)
}

// Get an actual client object from the Factory, from settings
func (clientFactory *FSouza_ClientFactory) MakeClient(logger log.Log, settings ClientSettings) (Client, bool) {
	client := &FSouza_Client{backend: clientFactory.client}
	return Client(client), client.Init(logger, clientFactory.conf, settings)
}

/**
 * Coach: Client
 */

// FSouza client settings struct
type FSouza_ClientSettings struct {
	log          log.Log
	conf         *conf.Project
	dependencies *Dependencies

	Author     string `json:"Author,omitempty" yaml:"Author,omitempty"`
	Repository string `json:"Repo,omitempty" yaml:"Repo,omitempty"`

	BuildPath string `json:"Build,omitempty" yaml:"Build,omitempty"`

	Config docker.Config     `json:"Config,omitempty" yaml:"Config,omitempty"`
	Host   docker.HostConfig `json:"Host,omitempty" yaml:"Host,omitempty"`
}

func (settings *FSouza_ClientSettings) Init(logger log.Log, project *conf.Project) bool {
	settings.log = logger
	settings.conf = project
	settings.dependencies = &Dependencies{}
	return true
}
func (settings *FSouza_ClientSettings) Prepare(logger log.Log, nodes *Nodes) bool {
	settings.log = logger

	/**
	 * Remap all binds so that they are based on the project folder
	 * - any relative path is remapped to be from the project root
	 * - any path starting with ~ is remapped to the user home flder
	 * - any absolute paths are left as is
	 */
	if settings.Host.Binds != nil {
		var binds []string
		for index, bind := range settings.Host.Binds {
			binds = strings.SplitN(bind, ":", 3)
			if path.IsAbs(binds[0]) {
				continue
			} else if binds[0][0:1] == "~" { // one would think that path can handle such terminology
				if rootPath, ok := settings.conf.Paths.Path("user-home"); ok {
					binds[0] = path.Join(rootPath, binds[0][1:])
				}
			} else {
				if rootPath, ok := settings.conf.Paths.Path("project-root"); ok {
					binds[0] = path.Join(rootPath, binds[0][:])
				}
			}
			settings.Host.Binds[index] = strings.Join(binds, ":")
		}
	}

	// build dependencies by looking at the Host Links and VolumesFrom lists
	settings.dependenciesFromConfig(logger, nodes, settings.Host.Links)
	settings.dependenciesFromConfig(logger, nodes, settings.Host.VolumesFrom)

	return true
}

// Add any nodes dependencies found when analyzing some of the various string slice docker configurations
func (settings *FSouza_ClientSettings) dependenciesFromConfig(logger log.Log, nodes *Nodes, config []string) {
	if config != nil {
		for _, item := range config {
			name := strings.SplitN(item, ":", 3)[0]
			if node, ok := nodes.Node(name); ok {
				settings.dependencies.SetDependency(name, Dependency(&NodeDependency{Node: node}))
			}
		}
	}
}

func (settings *FSouza_ClientSettings) Settings() interface{} {
	return settings
}
func (settings *FSouza_ClientSettings) nodeSettings(client *FSouza_Client, node Node) FSouza_ClientSettings {
	tokens := conf.Tokens{}
	tokens.SetToken("NODE", node.Id())
	tokens.SetToken("NODEMACHINE", node.MachineName())
	copy := settings.copy(tokens)

	return copy
}
func (settings *FSouza_ClientSettings) instancesSettings(client *FSouza_Client, instances Instances) FSouza_ClientSettings {
	return settings.copy(nil)
}
func (settings *FSouza_ClientSettings) instanceSettings(client *FSouza_Client, instance Instance) FSouza_ClientSettings {
	tokens := conf.Tokens{}
	tokens.SetToken("INSTANCE", instance.Id())
	tokens.SetToken("INSTANCEMACHINE", instance.MachineName())
	copy := settings.copy(tokens)

	if settings.Host.Links != nil && len(settings.Host.Links) > 0 {
		newLinks := []string{}
		for _, link := range settings.Host.Links {
			linkSplit := strings.SplitN(link, ":", 2)
			if transformedSet, found := settings.dependencies.DependencyIdTranform(linkSplit[0]); found {
				for _, transformed := range transformedSet {
					linkSplit[0] = transformed
					linkSplit[1] = strings.Replace(linkSplit[1], "%SOURCE", transformed, -1)
					newLinks = append(newLinks, strings.Join(linkSplit, ":"))
				}
			} else {
				newLinks = append(newLinks, link)
			}
		}
		copy.Host.Links = newLinks
	}
	if settings.Host.VolumesFrom != nil && len(settings.Host.VolumesFrom) > 0 {
		newVolumesFrom := []string{}
		for _, volumeFrom := range settings.Host.VolumesFrom {
			volumeFromSplit := strings.SplitN(volumeFrom, ":", 2)
			if transformedSet, found := settings.dependencies.DependencyIdTranform(volumeFromSplit[0]); found {
				for _, transformed := range transformedSet {
					volumeFromSplit[0] = transformed
					newVolumesFrom = append(newVolumesFrom, strings.Join(volumeFromSplit, ":"))
				}
			} else {
				newVolumesFrom = append(newVolumesFrom, volumeFrom)
			}
		}
		copy.Host.VolumesFrom = newVolumesFrom
	}

	return copy
}

func (settings *FSouza_ClientSettings) copy(tokens conf.Tokens) FSouza_ClientSettings {
	settings_json := settings.toJson()
	if tokens != nil {
		settings_json = tokens.TokenReplace(settings_json)
	}

	copy := FSouza_ClientSettings{}
	copy.fromJson(settings_json)
	return copy
}
func (settings *FSouza_ClientSettings) toJson() string {
	settings_json, _ := json.Marshal(settings)
	return string(settings_json)
}
func (settings *FSouza_ClientSettings) fromJson(source string) {
	json.Unmarshal([]byte(source), settings)
}

// FSouza Coach Client object
type FSouza_Client struct {
	settings FSouza_ClientSettings
	log      log.Log
	conf     *conf.Project
	backend  *FSouza_Wrapper

	id string
}

type FSouza_NodeClient struct {
	*FSouza_Client
	settings FSouza_ClientSettings
	node     Node
}

func (nodeClient *FSouza_NodeClient) Init(client *FSouza_Client, node Node) {
	nodeClient.FSouza_Client = client
	nodeClient.node = node
	nodeClient.settings = client.settings.nodeSettings(client, node)
}

type FSouza_InstancesClient struct {
	*FSouza_Client
	settings  FSouza_ClientSettings
	instances Instances
}

func (instancesClient *FSouza_InstancesClient) Init(client *FSouza_Client, instances Instances) {
	instancesClient.FSouza_Client = client
	instancesClient.instances = instances
	instancesClient.settings = client.settings.instancesSettings(client, instances)
}

type FSouza_InstanceClient struct {
	*FSouza_Client
	settings FSouza_ClientSettings
	instance Instance
}

func (instanceClient *FSouza_InstanceClient) Init(client *FSouza_Client, instance Instance) {
	instanceClient.FSouza_Client = client
	instanceClient.instance = instance
	instanceClient.settings = client.settings.instanceSettings(client, instance)
}

func (client *FSouza_Client) Init(logger log.Log, project *conf.Project, settings ClientSettings) bool {
	client.log = logger
	client.conf = project

	// make sure that the settings that were given, where the proper "FSouza_Client" type
	settingsTyped := settings.Settings()
	switch asserted := settingsTyped.(type) {
	case *FSouza_ClientSettings:
		client.settings = *asserted
		client.settings.Init(logger, project)
		return true
	default:
		logger.Error("Invalid settings type passed to Fsouza Client")
		return false
	}
}
func (client *FSouza_Client) Prepare(logger log.Log, nodes *Nodes, node Node) bool {
	client.id = node.MachineName()

	// the settings object has it's own Prepare
	client.settings.Prepare(logger, nodes)

	return true
}

func (client *FSouza_Client) NodeClient(node Node) NodeClient {
	nodeClient := &FSouza_NodeClient{}
	nodeClient.Init(client, node)
	return NodeClient(nodeClient)
}
func (client *FSouza_Client) InstancesClient(instances Instances) InstancesClient {
	instancesClient := &FSouza_InstancesClient{}
	instancesClient.Init(client, instances)
	return InstancesClient(instancesClient)
}
func (client *FSouza_Client) InstanceClient(instance Instance) InstanceClient {
	instanceClient := &FSouza_InstanceClient{}
	instanceClient.Init(client, instance)
	return InstanceClient(instanceClient)
}

func (client *FSouza_Client) Can(action string) bool {
	switch action {
	case "build":
		return client.settings.BuildPath != ""
	case "pull":
		return client.settings.BuildPath == "" && client.settings.Config.Image != ""
	default:
		return true
	}
}

func (client *FSouza_Client) DependsOn(target string) bool {
	_, ok := client.settings.dependencies.Dependency(target)

	return ok
}

/**
 * FSouza Client Wrapper
 */

type FSouza_Wrapper struct {
	*docker.Client

	cachedImages     []docker.APIImages
	cachedContainers []docker.APIContainers
}

// Init constructor for the client wrapper
func (wrapper *FSouza_Wrapper) Init(logger log.Log, settings FSouza_ClientFactorySettings) bool {
	var client *docker.Client
	var err error
	logger.Debug(log.VERBOSITY_DEBUG_WOAH, "Docker client conf: ", settings)

	if strings.HasPrefix(settings.Host, "tcp://") {

		if _, err := os.Stat(settings.CertPath); err == nil {
			// TCP DOCKER CLIENT WITH CERTS
			client, err = docker.NewTLSClient(
				settings.Host,
				path.Join(settings.CertPath, "cert.pem"),
				path.Join(settings.CertPath, "key.pem"),
				path.Join(settings.CertPath, "ca.pem"),
			)
		} else {
			// TCP DOCKER CLIENT WITHOUT CERTS
			client, err = docker.NewClient(settings.Host)
		}

	} else if strings.HasPrefix(settings.Host, "unix://") {
		// TCP DOCKER CLIENT WITHOUT CERTS
		client, err = docker.NewClient(settings.Host)
	} else {
		err = errors.New("Unknown client host :" + settings.Host)
	}

	if err == nil {
		logger.Debug(log.VERBOSITY_DEBUG_WOAH, "FSouza Docker client created:", client)
		wrapper.Client = client
		return true
	} else {
		logger.Error(err.Error())
		return false
	}
}

// Reload all of the client images and/or containers from the remote client
func (wrapper *FSouza_Wrapper) Refresh(refreshImages bool, refreshContainers bool) error {
	var err error

	if refreshImages {
		filters := map[string][]string{}
		//filters["id"] = []string{"something"} // As it turns out, filters are actually kind of useless : https://stackoverflow.com/questions/24659300/how-to-use-docker-images-filter

		options := docker.ListImagesOptions{
			Filters: filters,
		}
		wrapper.cachedImages, err = wrapper.ListImages(options)
	}
	if refreshContainers {
		filters := map[string][]string{}
		//filters["id"] = []string{"project"} // As it turns out, filters are actually kind of useless : https://stackoverflow.com/questions/24659300/how-to-use-docker-images-filter
		//filters["status"] = []string{"running"}

		options := docker.ListContainersOptions{
			All:     true,
			Filters: filters,
		}
		wrapper.cachedContainers, err = wrapper.ListContainers(options)
	}

	return err
}

// Return a list of all of the images registered on the client
func (wrapper *FSouza_Wrapper) AllImages(refresh bool) ([]docker.APIImages, error) {
	var err error
	if refresh || wrapper.cachedImages == nil {
		err = wrapper.Refresh(true, false)
	}
	return wrapper.cachedImages, err
}

// Return a list of remote images that have a specific string prefix repo tag
func (wrapper *FSouza_Wrapper) MatchImages(prefix string) ([]docker.APIImages, error) {
	prefix = strings.ToLower(prefix)
	images, err := wrapper.AllImages(false)
	filteredImages := []docker.APIImages{}
	for _, image := range images {
	eachimage:
		for _, tag := range image.RepoTags {
			if strings.HasPrefix(tag, prefix) {
				filteredImages = append(filteredImages, image)
				continue eachimage
			}
		}
	}
	return filteredImages, err
}

// Return a list of all of the containers registered on the client
func (wrapper *FSouza_Wrapper) AllContainers(refresh bool) ([]docker.APIContainers, error) {
	var err error
	if refresh || wrapper.cachedContainers == nil {
		err = wrapper.Refresh(false, true)
	}
	return wrapper.cachedContainers, err
}
func (wrapper *FSouza_Wrapper) MatchContainers(prefix string, running bool) ([]docker.APIContainers, error) {
	containers, err := wrapper.AllContainers(false)
	filteredContainers := []docker.APIContainers{}
	for _, container := range containers {
		if running && !strings.Contains(container.Status, "Up") {
			continue
		}

	eachcontainer:

		for _, containerName := range container.Names {
			if strings.Contains(containerName, "/"+prefix) {
				filteredContainers = append(filteredContainers, container)
				continue eachcontainer
			}
		}
	}
	return filteredContainers, err
}

/**
 * NodeClient meta-methods
 */

func (client *FSouza_Client) GetImageName() (image, tag string) {
	image = strings.ToLower(client.id)
	tag = "latest"

	if client.settings.Config.Image == "" {
		return
	}
	if strings.Contains(client.settings.Config.Image, ":") {
		split := strings.SplitN(client.settings.Config.Image, ":", 2)
		image = split[0]
		tag = split[1]
	} else {
		image = client.settings.Config.Image
	}
	return
}

func (client *FSouza_NodeClient) Images() []docker.APIImages {
	matchString := client.node.MachineName()
	if client.settings.Config.Image != "" {
		matchString = client.settings.Config.Image
	}

	images, _ := client.backend.MatchImages(matchString)
	return images
}
func (client *FSouza_NodeClient) HasImage() bool {
	return len(client.Images()) > 0
}

func (client *FSouza_NodeClient) NodeInfo(logger log.Log) {
	images := client.Images()

	if len(images) == 0 {
		client.log.Message("|-- no image [" + client.node.MachineName() + "]")
	} else {
		client.log.Message("|-> Images")

		w := new(tabwriter.Writer)
		w.Init(client.log, 8, 12, 2, ' ', 0)

		row := []string{
			"|=",
			"ID",
			"RepoTags",
			"Created",
			//		"Size",
			//		"VirtualSize",
			//	"ParentID",
			//		"RepoDigests",
			//		"Labels",
		}
		w.Write([]byte(strings.Join(row, "\t") + "\n"))

		for _, image := range images {
			row := []string{
				"|-",
				image.ID[:11],
				strings.Join(image.RepoTags, ","),
				strconv.FormatInt(image.Created, 10),
				//			strconv.FormatInt(image.Size, 10),
				//			strconv.FormatInt(image.VirtualSize, 10),
				//			image.ParentID,
				// 			strings.Join(image.RepoDigests, "\n"),
				// 			strings.Join(image.Labels, "\n"),
			}
			w.Write([]byte(strings.Join(row, "\t") + "\n"))
		}
		w.Flush()
	}
}

/**
 * InstancesInfo Interface
 */

func (client *FSouza_InstancesClient) InstancesInfo(logger log.Log) {
	instances := client.instances

	if instances.MachineName() == INSTANCES_NULL_MACHINENAME {

	} else if len(instances.InstancesOrder()) == 0 {
		logger.Message("|-= no containers")
	} else {
		logger.Message("|-> instances (containers) MachineName:" + instances.MachineName())

		w := new(tabwriter.Writer)
		w.Init(logger, 8, 12, 2, ' ', 0)

		row := []string{
			"|=",
			"Name",
			"Container",
			"Default",
			"Created",
			"Running",
			"Status",
			"ID",
			"Created",
			"Names",
		}
		w.Write([]byte(strings.Join(row, "\t") + "\n"))

		for _, name := range instances.InstancesOrder() {
			instance, _ := instances.Instance(name)
			machineName := instance.MachineName()
			instanceClient := instance.Client()

			row := []string{
				"|-",
				name,
				machineName,
			}
			if instance.IsDefault() {
				row = append(row, "yes")
			} else {
				row = append(row, "no")
			}
			if instanceClient.HasContainer() {
				row = append(row, "yes")
			} else {
				row = append(row, "no")
			}
			if instanceClient.IsRunning() {
				row = append(row, "yes")
			} else {
				row = append(row, "no")
			}

			containers, _ := client.backend.MatchContainers(machineName, false)
			for _, container := range containers {
				row = append(row,
					container.Status,
					container.ID[:12],
					strconv.FormatInt(int64(container.Created), 10),
					strings.Join(container.Names, ", "),
				)
				break
			}

			w.Write([]byte(strings.Join(row, "\t") + "\n"))
		}
		w.Flush()

	}
}

func (client *FSouza_InstancesClient) InstancesFound(logger log.Log) []string {
	ids := []string{}
	return ids
}

/**
 * InstanceClient meta-methods
 */

func (client *FSouza_InstancesClient) Containers(running bool) []docker.APIContainers {
	instances := client.instances
	matchString := instances.MachineName()
	if matchString == INSTANCES_NULL_MACHINENAME {
		return []docker.APIContainers{}
	} else {
		containers, _ := client.backend.MatchContainers(matchString, running)
		return containers
	}
}
func (client *FSouza_InstanceClient) Containers(running bool) []docker.APIContainers {
	instance := client.instance
	matchString := instance.MachineName()
	if matchString == INSTANCES_NULL_MACHINENAME {
		return []docker.APIContainers{}
	} else {
		containers, _ := client.backend.MatchContainers(matchString, running)
		return containers
	}
}
func (client *FSouza_InstanceClient) HasContainer() bool {
	return len(client.Containers(false)) > 0
}
func (client *FSouza_InstanceClient) IsRunning() bool {
	return len(client.Containers(true)) > 0
}

/**
 * NodeClient interface: Operation Methods
 */

func (client *FSouza_NodeClient) Build(logger log.Log, force bool) bool {
	image, tag := client.GetImageName()

	if client.settings.BuildPath == "" {
		logger.Warning("Node image [" + image + ":" + tag + "] not built as an empty path was provided.  You must point Build: to a path inside .coach")
		return false
	}

	if !force && client.HasImage() {
		logger.Warning("Node image [" + image + ":" + tag + "] not built as an image already exists.  You can force this operation to build this image")
		return false
	}

	// determine an absolute buildPath to the build, for Docker to use.
	buildPath := ""
	for _, confBuildPath := range client.conf.Paths.GetConfSubPaths(client.settings.BuildPath) {
		logger.Debug(log.VERBOSITY_DEBUG_STAAAP, "Looking for Build: "+confBuildPath)
		if _, err := os.Stat(confBuildPath); !os.IsNotExist(err) {
			buildPath = confBuildPath
			break
		}
	}
	if buildPath == "" {
		logger.Error("No matching build path could be found [" + client.settings.BuildPath + "]")
	}

	options := docker.BuildImageOptions{
		Name:           image + ":" + tag,
		ContextDir:     buildPath,
		RmTmpContainer: true,
		OutputStream:   logger,
	}

	logger.Info("Building node image [" + image + ":" + tag + "] From build path [" + buildPath + "]")

	// ask the docker client to build the image
	err := client.backend.BuildImage(options)

	if err != nil {
		logger.Error("Node build failed [" + client.node.MachineName() + "] in build path [" + buildPath + "] => " + err.Error())
		return false
	} else {
		logger.Message("Node succesfully built image [" + image + ":" + tag + "] From path [" + buildPath + "]")
		return true
	}

}

func (client *FSouza_NodeClient) Destroy(logger log.Log, force bool) bool {
	// Get the image name
	image, tag := client.GetImageName()
	if tag != "" {
		image += ":" + tag
	}

	if !client.HasImage() {
		logger.Warning("Node has no image to destroy [" + image + "]")
		return false
	}

	options := docker.RemoveImageOptions{
		Force: force,
	}

	// ask the docker client to remove the image
	err := client.backend.RemoveImageExtended(image, options)

	if err != nil {
		logger.Error("Node image removal failed [" + image + "] => " + err.Error())
		return false
	} else {
		logger.Message("Node image was removed [" + image + "]")
		return true
	}
}

func (client *FSouza_NodeClient) Pull(logger log.Log, force bool) bool {
	image, tag := client.GetImageName()
	actionCacheTag := "pull:" + image + ":" + tag

	if _, ok := actionCache[actionCacheTag]; ok {
		logger.Message("Node image [" + image + ":" + tag + "] was just pulled, so not pulling it again.")
		return true
	}

	if !force && client.HasImage() {
		logger.Info("Node already has an image [" + image + ":" + tag + "], so not pulling it again.  You can force this operation if you want to pull this image.")
		return false
	}

	options := docker.PullImageOptions{
		Repository:    image,
		OutputStream:  logger,
		RawJSONStream: false,
	}

	if tag != "" {
		options.Tag = tag
	}

	var auth docker.AuthConfiguration
	// 		var ok bool
	//options.Registry = "https://index.docker.io/v1/"

	// 		auths, _ := docker.NewAuthConfigurationsFromDockerCfg()
	// 		if auth, ok = auths.Configs[registry]; ok {
	// 			options.Registry = registry
	// 		} else {
	// 			node.log.Warning("You have no local login credentials for any repo. Defaulting to no login.")
	auth = docker.AuthConfiguration{}
	options.Registry = "https://index.docker.io/v1/"
	// 		}

	logger.Message("Pulling node image [" + image + ":" + tag + "] from server [" + options.Registry + "] using auth [" + auth.Username + "] : " + image + ":" + tag)
	logger.Debug(log.VERBOSITY_DEBUG_LOTS, "AUTH USED: ", map[string]string{"Username": auth.Username, "Password": auth.Password, "Email": auth.Email, "ServerAdddress": auth.ServerAddress})

	// ask the docker client to build the image
	err := client.backend.PullImage(options, auth)

	if err != nil {
		logger.Error("Node image not pulled : " + image + " => " + err.Error())
		actionCache[actionCacheTag] = false
		return false
	} else {
		logger.Message("Node image pulled: " + image + ":" + tag)
		actionCache[actionCacheTag] = false
		return true
	}
}

/**
 * InstanceClient : Action methods
 */

func (client *FSouza_InstanceClient) Attach(logger log.Log) bool {
	id := client.instance.MachineName()

	// build options for the docker attach operation
	options := docker.AttachToContainerOptions{
		Container:    id,
		InputStream:  os.Stdin,
		OutputStream: os.Stdout,
		ErrorStream:  logger,

		Logs:   true, // Get container logs, sending it to OutputStream.
		Stream: true, // Stream the response?

		Stdin:  true, // Attach to stdin, and use InputStream.
		Stdout: true, // Attach to stdout, and use OutputStream.
		Stderr: true,

		//Success chan struct{}

		RawTerminal: client.settings.Config.Tty, // Use raw terminal? Usually true when the container contains a TTY.
	}

	logger.Message("Attaching to instance container [" + id + "]")
	err := client.backend.AttachToContainer(options)
	if err != nil {
		logger.Error("Failed to attach to instance container [" + id + "] =>" + err.Error())
		return false
	} else {
		logger.Message("Disconnected from instance container [" + id + "]")
		return true
	}
}

func (client *FSouza_InstanceClient) Create(logger log.Log, overrideCmd []string, force bool) bool {
	instance := client.instance

	if !force && client.HasContainer() {
		logger.Info("[" + instance.MachineName() + "]: Skipping node instance, which already has a container")
		return false
	}

	/**
	* Transform node data, into a format that can be used
	* for the actual Docker call.  This involves transforming
	* the node keys into docker container ids, for things like
	* the name, Links, VolumesFrom etc
	 */
	name := instance.MachineName()
	Config := client.settings.Config
	Host := client.settings.Host

	image, tag := client.GetImageName()
	if tag != "" && tag != "latest" {
		image += ":" + tag
	}
	Config.Image = image

	if len(overrideCmd) > 0 {
		Config.Cmd = overrideCmd
	}

	// ask the docker client to create a container for this instance
	options := docker.CreateContainerOptions{
		Name:       name,
		Config:     &Config,
		HostConfig: &Host,
	}

	container, err := client.backend.CreateContainer(options)
	client.backend.Refresh(false, true)

	if err != nil {

		logger.Debug(log.VERBOSITY_DEBUG, "CREATE FAIL CONTAINERS: ", err)

		/**
		* There is a weird bug with the library, where sometimes it
		* reports a missing image error, and yet it still creates the
		* container.  It is not clear if this failure occurs in the
		* remote API, or in the dockerclient library.
		 */

		if err.Error() == "no such image" && client.HasContainer() {
			logger.Message("Created instance container [" + name + " FROM " + Config.Image + "] => " + container.ID[:12])
			logger.Warning("Docker created the container, but reported an error due to a 'missing image'.  This is a known bug, that can be ignored")
			return true
		}

		logger.Error("Failed to create instance container [" + name + " FROM " + Config.Image + "] => " + err.Error())
		return false
	} else {
		logger.Message("Created instance container [" + name + "] => " + container.ID[:12])
		return true
	}
}

func (client *FSouza_InstanceClient) Remove(logger log.Log, force bool) bool {
	name := client.instance.MachineName()
	options := docker.RemoveContainerOptions{
		ID: name,
	}

	// ask the docker client to remove the instance container
	err := client.backend.RemoveContainer(options)

	if err != nil {
		logger.Error("Failed to remove instance container [" + name + "] =>" + err.Error())
		return false
	} else {
		logger.Message("Removed instance container [" + name + "] ")
		return true
	}

	return false
}

func (client *FSouza_InstanceClient) Start(logger log.Log, force bool) bool {
	// Convert the node data into docker data (transform node keys to container IDs for things like Links & VolumesFrom)
	id := client.instance.MachineName()
	Host := client.settings.Host

	// ask the docker client to start the instance container
	err := client.backend.StartContainer(id, &Host)

	if err != nil {
		logger.Error("Failed to start node container [" + id + "] => " + err.Error())
		return false
	} else {
		logger.Message("Node instance started [" + id + "]")
		return true
	}
}

func (client *FSouza_InstanceClient) Stop(logger log.Log, force bool, timeout uint) bool {
	id := client.instance.MachineName()

	err := client.backend.StopContainer(id, timeout)
	if err != nil {
		logger.Error("Failed to stop node container [" + id + "] => " + err.Error())
		return false
	} else {
		logger.Message("Node instance stopped [" + id + "]")
		return true
	}
}

func (client *FSouza_InstanceClient) Pause(logger log.Log) bool {
	id := client.instance.MachineName()

	err := client.backend.PauseContainer(id)
	if err != nil {
		logger.Error("Failed to pause intance [" + client.instance.Id() + "] Container [" + id + "] =>" + err.Error())
		return false
	} else {
		logger.Message("Paused instance [" + client.instance.Id() + "] Container [" + id + "]")
		return true
	}
}

func (client *FSouza_InstanceClient) Unpause(logger log.Log) bool {
	id := client.instance.MachineName()

	err := client.backend.UnpauseContainer(id)
	if err != nil {
		logger.Error("Failed to unpause Instance [" + client.instance.Id() + "] Container [" + id + "] =>" + err.Error())
		return false
	} else {
		logger.Message("Unpaused Instance [" + client.instance.Id() + "] Container [" + id + "]")
		return true
	}
}

func (client *FSouza_InstanceClient) Commit(logger log.Log, tag string, message string) bool {
	id := client.instance.MachineName()
	config := client.settings.Config
	repo := client.settings.Repository
	author := client.settings.Author

	if repo == "" {
		repo, _ = client.GetImageName()
	}

	options := docker.CommitContainerOptions{
		Container:  id,
		Repository: repo,
		Tag:        tag,
		Run:        &config,
	}

	if message != "" {
		options.Message = message
	}
	if author != "" {
		author = client.conf.Author
	}

	_, err := client.backend.CommitContainer(options)
	if err != nil {
		logger.Warning("Failed to commit container changes to an image [" + client.instance.Id() + ":" + id + "] : " + tag)
		return false
	} else {
		logger.Message("Committed container changes to an image [" + client.instance.Id() + ":" + id + "] : " + tag)
		return true
	}
}

func (client *FSouza_InstanceClient) Run(logger log.Log, persistant bool, cmd []string) bool {
	instance := client.instance

	// Set up some additional settings for TTY commands
	if client.settings.Config.Tty == true {

		// set a default hostname to make a prettier prompt
		if client.settings.Config.Hostname == "" {
			client.settings.Config.Hostname = instance.Id()
		}

		// make sure that all tty runs have openstdin
		client.settings.Config.OpenStdin = true
	}

	client.settings.Config.AttachStdin = true
	client.settings.Config.AttachStdout = true
	client.settings.Config.AttachStderr = true

	// 1. get the container for the instance (create it if needed)
	hasContainer := client.HasContainer()
	if !hasContainer {
		logger.Info("Creating new disposable RUN container")

		if hasContainer = client.Create(logger, cmd, false); hasContainer {
			logger.Debug(log.VERBOSITY_DEBUG, "Created disposable run container")
			if !persistant {
				// 5. [DEFERED] remove the container (if not instructed to keep it)
				defer client.Remove(logger, true)
			}
		} else {
			logger.Error("Failed to create disposable run container")
		}
	} else {
		logger.Info("Run container already exists")
	}

	if hasContainer {

		// 3. start the container (set up a remove)
		logger.Info("Starting RUN container")
		ok := client.Start(logger, false)

		// 4. attach to the container
		if ok {
			logger.Info("Attaching to disposable RUN container")
			client.Attach(logger)
			return true
		} else {
			logger.Error("Could not start RUN container")
			return false
		}

	} else {
		logger.Error("Could not create RUN container")
	}
	return false
}
