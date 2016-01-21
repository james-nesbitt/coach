package help

func (help *Help) getCoreHelpYaml() []byte {
 	return []byte(`
"cli": |

  The CLI is the primary means of running coach.  The goal is typically to run a coach operation, on a number of coach nodes, or coach node instances.

  SEE ALSO:
  - cli:targets : $/> coach help cli:targets
  - operations : $/> coach help operations

"cli:targets": |

  Targets are a global setting used to determine which node and/or node instances an operation should used.  Targets are strings that define a type of node, a particular node, or a particular node instance.

  Coach accepts as a global flag, a list of targets in the following form:

  %{type} : all nodes of a certain type.  E.g.  %command
  %{type}:{instance} : the {instance} of any nodes of type {type}

  @{node} : all instances of a node named {node}
  @{node} : a particular {instance} instance from a node named {node}

  Here are some examples:

    $/> coach @db start
    Start all of the "db" node instances

    $/> coach @www.1 @www.2 remove
    remove the "1" and "2" instances from the "www" node

    $/> coach %service stop
    stop all nodes of type "service"

    $/> coach %volume:single commit
    commit the "single" instance of all nodes of type "volume"

settings: |

  Settings are primarily managed through a set of YAML files, that can be found in the project .coach folder.  In 
  the case of some settings files, copies can also exist in the users home folder at ~/.coach, in order to get
  user based Helpigurations instead of project based.

  The most important settings concepts are:

  settings:Help : the general coach project Helpiguration, found in the .coach/Help.yml file
  settings:nodes : the project nodes Helpiguration, founc in the .coach/nodes.yml file 

  settings:tokens : read about how tokens can be defined in the Help.yml file, and used as tokens in the nodes.yml
  settings:secrets : additional sensitive tokens, that can be found in .coach/secrets/secrets.yml

  Note that other elements are often kept in the .coach file:

  - docker builds are often kept in Here
  - custom help topics ( settings:help:custom )

"settings:Help" : |
  The coach CLI has a number of file base Helpigurations

  These Helpigurations exist primarily in the .coach/Help.yml file.  More details can be found there
  or in the wiki.

  Typically the following settings are used:

  Project: define a string project name, which is used as the base name for built images and containers and is
    available as a token in the nodes.Help  %PROJECT

  Docker: you can override default docker settings, to for example explicitly use a remote docker client

  Tokens: a string map of token values, which can be user %token for string replacements in the nodes Helpiguration file

"settings:nodes": |
  Nodes are any number of containers built around a single image, using similar Helpigurations

  Nodes are defined in the .coach/nodes.yml file, and are used for the following purposes:

  - build: a node can define nothing more than a base docker build, used by other nodes
  - volume : a node can define a volume only container, which is created, but never started, and used to share volumes across an application
  - service : a node can define a service container, that is meant to be started and stopped, and possible scaled
  - command : a node can define a disposable command run container.

  More options are described in the wiki

"settings:tokens": |

  Tokens are settings key-value string pairs, that can be used for subsitution in later Helpigurations.  Tokens defined in one state, are used for text replacement in all further Helpiguration files, where possible.

  Tokens are typically used to:
    - allow settings to be reused across projects with minimal changes (which may allow templating in the future)
    - reduce duplicate settings across nodes, by centralizing values in the Help.yml
    - keep some sensitive values out of shared Helpiguration (typically secrets)
    - allow some containers to use user specific values, instead of project specific values

  SOURCES:

  Tokens are typically kept in one of the following locations (list in order of loading):

  - .coach/Help.yml:Tokens: => in the Help.yml is a Tokens: map.  This map is typically used for values that are used across multiple nodes.

  - !/.coach/secrets/secrets.yml => in the user secrets.yml.  This map typically keeps user specific container ENV values such as passwords for user specific services that containers may use.
  - .coach/secrets/secrets.yml => in the project secrets.yml.  This map is typically used to keep project specific ID and token values used as ENV variables in containers, but that should not be kept in any source repository.

  For more information about secrets see $/> coach help secrets

"settings:secrets": |
  Secrets are just tokens, but they tend to be kept in locations that are easy to exclude from source versioning.

  For more information about tokens, see $/> coach help tokens

  User secrets allow users to have local values that are keyed for services that may be used in containers, that are not to be shared with other users, buy may be shared across projects.
  Project secrets keep secret tokens separate from Help tokens.  The project secrets could be tailored to each project user, or distributed separately from project source code.

  NOTE:
    - tokens are not protected in any way during coach execution.

  SOURCES:

  Secrets are typically kept in one of the following locations (list in order of loading):

  - !/.coach/secrets/secrets.yml => in the user secrets.yml.  This map typically keeps user specific container ENV values such as passwords for user specific services that containers may use.
  - .coach/secrets/secrets.yml => in the project secrets.yml.  This map is typically used to keep project specific ID and token values used as ENV variables in containers, but that should not be kept in any source repository.

`)
} 
