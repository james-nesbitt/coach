package libs


type BuildNodeSettings struct {
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
}

type BuildNode struct {
	Settings BuildNodeSettings
	BaseNode
}
