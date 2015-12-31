package libs

type CommandNodeSettings struct {
	Type string `json:"Type,omitempty" yaml:"Type,omitempty"`
}

type CommandNode struct {
	Settings CommandNodeSettings
	BaseNode
}
