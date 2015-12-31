package libs

type Targets []*Target

type Target struct {
	Node      *Node
	Instances []*Instance
}

func (targets *Targets) fromNodes(identifiers []string, nodes Nodes) Targets {
	return nil
}
