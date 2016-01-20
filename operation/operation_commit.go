package operation

import (
	"strings"

	"github.com/james-nesbitt/coach-tools/libs"
	"github.com/james-nesbitt/coach-tools/log"
)

type CommitOperation struct {
	log     log.Log
	targets *libs.Targets

	force bool

	tag     string
	message string
}

func (operation *CommitOperation) Id() string {
	return "commit"
}
func (operation *CommitOperation) Flags(flags []string) bool {
	for index := 0; index < len(flags); index++ {
		flag := flags[index]

		switch flag {
		case "-t":
			fallthrough
		case "--tag":
			if !strings.HasPrefix(flags[index+1], "-") {
				index++
				operation.tag = flags[index]
			}
		case "-m":
			fallthrough
		case "--message":
			if !strings.HasPrefix(flags[index+1], "-") {
				index++
				operation.message = flags[index]
			}
		}
	}
	return true
}

func (operation *CommitOperation) Help(topics []string) {
	operation.log.Message(`Operation: Commit

Coach will attempt to commit a container to it's image.

SYNTAX:
	$/> coach {targets} commit [--tag {tag}] [--repo {repo}] [--message "{message}"]

	{targets} what target node instances the operation should process ($/> coach help targets)
	--tag "{tag}" : what image tag to use (default: "latest")
	--message "{message}" : what commit message to use

ACCESS:
	- only nodes with the "commit" access are processed.  This excludes build nodes

`)
}
func (operation *CommitOperation) Run(logger log.Log) bool {
	logger.Info("Running operation: commit")
	logger.Debug(log.VERBOSITY_DEBUG, "Run:Targets", operation.targets.TargetOrder())

	for _, targetID := range operation.targets.TargetOrder() {
		target, targetExists := operation.targets.Target(targetID)
		if !targetExists {
			// this is strange
			logger.Warning("Internal target error, was told to use a target that doesn't exist")
			continue
		}

		node, hasNode := target.Node()
		instances, hasInstances := target.Instances()

		nodeLogger := logger.MakeChild(targetID)

		if !hasNode {
			nodeLogger.Warning("No node [" + node.MachineName() + "]")
		} else if !node.Can("Commit") {
			nodeLogger.Info("Node doesn't Commit [" + node.MachineName() + ":" + node.Type() + "]")
		} else if !hasInstances {
			nodeLogger.Info("No valid instances specified in target list [" + node.MachineName() + "]")
		} else {
			nodeLogger.Message("Creating instance containers")
			for _, id := range instances.InstancesOrder() {
				instance, _ := instances.Instance(id)
				instance.Client().Commit(logger, operation.tag, operation.message)
			}
		}
	}

	return true
}
