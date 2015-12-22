package libs


type Node interface {
	May(operation string) bool
	Client() *Client
	Instances(instance string) *Instances
}

type BaseNode struct {
	client *NodeClient
}
