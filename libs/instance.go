package libs

type Instances interface {
	Can(action string) bool
	Client() *Client
}

type BaseInstances struct {
	Settings *InstanceSettings
	Client *InstanceClient
}
func (instance *BaseInstance) Can(action string) bool {
	return true
}
func (instance *BaseInstance) Client() *Client {
	return instance.Client
}