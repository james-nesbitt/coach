package main

type Operation_Run struct {
	log Log

	Nodes Nodes
	Targets []string

	cmd []string
}
func (operation *Operation_Run) Flags(flags []string) {
	operation.cmd = flags
}
func (operation *Operation_Run) Run() {
	operation.log.Message("running run operation")
	operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.Targets)

// 	operation.Nodes.log = operation.log.ChildLog("OPERATION:RUN")
	operation.Nodes.Run(operation.Targets, operation.cmd)
}

func (nodes *Nodes) Run(targets []string, cmd []string) {
	for _, target := range nodes.GetTargets(targets) {
		target.log = nodes.log.ChildLog("NODE:"+target.Name)
		target.Run(cmd)
	}
}

func (node *Node) Run(cmd []string) bool {
	if node.Do("run") {

		id := "run" // perhaps we should randomly generate this to allow for persistance run containers

		// 1. create a new temporary instance for a node
		node.AddTemporaryInstance(id)
		instance := node.GetInstance(id)

		// 2. create the container for the new instance
		instance.Config.AttachStdin = true
		instance.Config.AttachStdout = true
		instance.Config.AttachStderr = true

		if instance.Create(cmd, false) {

		// 3. start the container (set up a remove)
			ok := instance.Start(false)

		// 4. attach to the container
			if ok {
				defer instance.Remove(true)
				instance.Attach()
				return true
			} else {
				return false
			}

		}

	}
	return false
}
