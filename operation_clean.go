package main


type Operation_Clean struct {
  log Log

  nodes Nodes
  targets []string

  force bool
  wipe bool
  timeout uint
}
func (operation *Operation_Clean) Flags(flags []string) {
  operation.force = false
  operation.wipe = false
  operation.timeout = 10

  for _, flag := range flags {
    switch flag {
      case "-f":
        fallthrough
      case "--force":
        operation.force = true
      case "-w":
        fallthrough
      case "--wipe":
        operation.wipe = true
      case "-q":
        fallthrough
      case "--quick":
        operation.timeout = 0
    }
  }
}

func (operation *Operation_Clean) Help(topics []string) {
  operation.log.Note(`Operation: Clean

Coach will attempt to Clean any node containers that are active.  Cleaning
means to stop and remove any active containers.  Wiping means to also remove
any built images.

Syntax:
    $/> coach {targets} Clean

    $/> coach {targets} Clean --wipe   
      also remove any built images

    - to eliminate any timeout delays use this Syntax
    $/> coach {targets} Clean --quick

  {targets} what target node instances the operation should process ($/> coach help targets)

`)
}

func (operation *Operation_Clean) Run() {
  operation.log.Info("running clean operation")
  operation.log.DebugObject(LOG_SEVERITY_DEBUG_LOTS, "Targets:", operation.targets)

//  operation.Nodes.log = operation.log.ChildLog("OPERATION:CLEAN")
  operation.nodes.Clean(operation.targets, operation.wipe, operation.force, operation.timeout)
}

// Clean a nodes list
func (nodes *Nodes) Clean(targets []string, wipe bool, force bool, timeout uint) {
  for _, target := range nodes.GetTargets(targets) {
    target.node.log.Message(target.node.Name+": Cleaning Node") 

    for _, instance := range target.instances {
      if instance.HasContainer(false) {
        instance.Clean(force, timeout)
      } else {
        instance.Node.log.Info(instance.Node.Name+": skipping node instance clean, as it has no container")
      }
    }

    if wipe {
      // trap any messages from the other intance operations
      if target.node.log.Severity()==LOG_SEVERITY_MESSAGE {
        target.node.log.Hush()
        defer target.node.log.UnHush()
      }

      target.node.Destroy(false)
      target.node.log.Info(target.node.Name+": node build cleaned") 
    }
  }
}

// Clean a node
func (node *Node) Clean(filters []string, wipe bool, force bool, timeout uint) {
  var instances []*Instance

  if len(filters)==0 {
    instances = node.GetInstances()
  } else {
    instances = node.FilterInstances(filters)
  }

  node.log.Message(node.Name+": Cleaning Node")
node.log.DebugObject(LOG_SEVERITY_ERROR, "NODE INSTANCES:", instances)
  for _, instance := range instances {
    if instance.HasContainer(false) {
      instance.Clean(force, timeout)
    } else {
      node.log.Info(node.Name+": skipping node instance clean, as it has no container")
    }
  }

  if wipe {
    // trap any messages from the other intance operations
    if node.log.Severity()==LOG_SEVERITY_MESSAGE {
      node.log.Hush()
      defer node.log.UnHush()
    }

    node.Destroy(force)
    node.log.Info(node.Name+": node build cleaned") 
  }

}

//Clean a node instance
func (instance *Instance) Clean(force bool, timeout uint) bool {
  instance.Node.log.Message(instance.Node.Name+": cleaning node instance ["+instance.Name+"]")

  if instance.HasContainer(true) {
    if !instance.Stop(force, timeout) {
      return false
    }
  }
  return instance.Remove(force)
}