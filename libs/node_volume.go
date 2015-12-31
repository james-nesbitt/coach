package libs

type VolumeNodeSettings struct {
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
}

type VolumeNode struct {
	Settings VolumeNodeSettings
	BaseNode
}
