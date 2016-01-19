package help

import (
	"github.com/james-nesbitt/coach-tools/log"
	"github.com/james-nesbitt/coach-tools/conf"
)

type Help struct {
	conf *conf.Project
	log log.Log

	topics map[string]string   `yaml:"Topics,omitempty"`
}

func (helper *Help) Init(logger log.Log, project *conf.Project) {
	helper.conf = project
	helper.log = logger
	helper.topics = map[string]string{}

  // add the core help
  helper.from_CoreHelpYaml(logger)

  // Add the Yaml file help
  helper.from_HelpYaml(logger, project)
}

func (helper *Help) merge(merge Help) {
	if helper.topics==nil {
		helper.topics = map[string]string{}
	}
  for key, topic := range merge.topics {
  	if _, ok := helper.topics[key]; !ok {
  		helper.topics[key] = topic
  	}
  }
}
func (helper *Help) Topic(name string, flags []string) (topic string, ok bool) {
	topic, ok = helper.topics[name]
	return
}
func (helper *Help) SetTopic(name string, topic string) bool {
	helper.topics[name] = topic
	return true
}
