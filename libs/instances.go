package libs

type Instances interface {

}

type BaseInstances struct {
	SettingsBase *InstanceSettings
	Instances map[string]*Instance
	Client *Client
}