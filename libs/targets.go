package libs

type Targets []*Target

type Target struct {
	Node      *Node
	Instances []*Instance
}
