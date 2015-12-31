package libs

type ServiceNodeSettings struct {
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
}

type ServiceNode struct {
	Settings ServiceNodeSettings
	BaseNode
}
