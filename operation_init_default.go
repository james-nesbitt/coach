package main

func (operation *Operation_Init) Init_Default_Run(source string, tasks *InitTasks) bool {

  if source=="" {
    source="bare"
  }

  switch source {
    case "bare":
      return operation.Init_Default_Bare(tasks)
    case "starter":
      return operation.Init_Default_Starter(tasks)
  }


  return false
}

func (operation *Operation_Init) Init_Default_Bare(tasks *InitTasks) bool {

  tasks.AddFile(".gitignore",`# Ignore coach secrets
.coach/secrets
`)
  tasks.AddFile(".coach/conf.yml",  `# Coach project conf
Project: bare`)
  tasks.AddFile(".coach/nodes.yml",  `# Project nodes
`)
  tasks.AddFile(".coach/secrets/secrets.yml", `# Private project tokens
`)
  tasks.AddFile("app/README.md", `# Bare Project
## /.coach/

  Project coach configuration

## /app

  Project source-code and assets path

`)

  tasks.AddMessage("Created local project as a `bare` project")

  return true
}

func (operation *Operation_Init) Init_Default_Starter(tasks *InitTasks) bool {

  tasks.AddFile(".gitignore",  `# Ignore coach secrets
.coach/secrets
`)
  tasks.AddFile(".coach/conf.yml",  `# Coach project conf
#
# Configurations for the coach project
#
# Project: project name, used to create container names
# Author: used for docker commits & pushes
#
# Tokens: a string map for string substitutions elsewhere
#    by using the format %{key} in files loaded after this
#    one
#
# Paths: a string map of paths that you can use for 
#    things like docker builds and mounts.  Paths can be 
#    absolute or relative to the project root.
#    Note that you can override the following paths:
#      - usertemplates : path to users init templates
#      - usersecrets : path to user secrets
#      - projectsecrets : path to project secrets
#      - build : path to builds
#    Paths are also made available as tokens using the
#    format PATH_{UPPERCASEKEY}
#

Project: starter

Author: me

Paths:
  key: "path"

Tokens:
  TOKEN_KEY: "TOKEN VALUE"

Settings:
  UseEnvVariablesAsTokens: "yes"   # include all of the user's ENV variables as possible tokens

Docker:  # Override Docker configuration
#  Host: "tcp://10.0.42.1"         # point to a remote docker server

`)
  tasks.AddFile(".coach/secrets/secrets.yml", `# Private project tokens
#
# Private tokens, which can be used in the nodes configuration, but
# can be kept out of the 
#
# {key}: {value}
#
# For substitution, use %{key} in other files
#
PERSONAL_APP_KEY: APP_KEY_VALUE
`)
  tasks.AddFile(".coach/nodes.yml",  `###
# Project nodes
#
# A string map of node configurations, each entry is a node, which 
# corresponds to a set of containers based on a single image.  the
# nodes have different types, an different purposes:
#
# Nodes: 
#   Each element in the yml file is a keyed configuration to configure
#   a single coach node, which is any number of docker containers
#   that use a single build/image.  The nodes can be just used to 
#   build images, or can be used to store data, run services or run
#   commands.  A node can use multiple containers to provide scaling
#   services, or alternate services.
#
# Types:
#   - build: only used for building images, let's you build a base 
#        image for other nodes
#   - volume: non-running containers meant to hold files that are 
#        shared with other nodes
#   - service: a node with containers that get started and stopped
#   - command: a node that uses disposable containers to run commands
#        inside an environment
#
# Instances:
#   Nodes can have multiple containers called instances. The instances
#   setting tells coach how to treat the node instances:
#   - scale: allow numeric instances, and allow the scale operation 
#        to spin up new instanes.
#   - fixed: if you provide a string list of instances, then they 
#         can be started and stopped by name.
#   - temporary: instances started are removed when stopped.  For example
#         command containers are considered temporary.
#
#   Multiple instances are complex, especially when it comes to linking
#   Nodes together.  If one node is linked to another, either by Link,
#   or VolumesFrom, then coach tries to link matching instance names
#   possible (or you can user node:all to link to all)
#
# Tokens:
#
#  Tokens can be used for string substitution in this file.  Tokens can come
#  from the conf.yml, or from the secrets.yml (or from the users secrets.)
#  Tokens are substituted using a "%" prefix, before the YML is parsed, so
#  you may need to wrap tokenized settings in "quotes" when using them.
#
#  The following tokens are available:
#
#    - anything from the tokens section of the conf.yml
#    - anything from the project or user secrets.yml
#    
#    - the project name is accessible using %PROJECT
#    - any path value is accessible using %PATH_UPPERCASEPATHKEY
#    - multi-instance nodes can access the intance name as %INSTANCE
#
# Coach configuration:
#
#   - Type: You can define the node type (or it will default to service)
#   - Build: You can assign a Build path, which will be used to build a local
#        image for the node.  If you also define an image then that
#        is used for the image name.  If the path is relative, then it
#        is considered relative inside the .coach/ folder.
#   - Instances: a string list of instance names, or the keyword temporary
#        or the keyword scaled, followed by the integer values for initial
#        and maximum number of running containers
#        e.g.:
#           Instances: first second third
#           Instances: temporary
#           Instances: scaled 3 9
#
# Docker remote API Configurations:
#
#   The two principle parts of the configuration for a node are the 
#   Config and Host settings.
#   These wrap the client config and can use any of it's keys, which 
#   gives access to pretty much all of the docker remote API settings.
#
#   - Config : the dockerclient config element used define the internals
#        of a container.  
#        https://github.com/fsouza/go-dockerclient/blob/master/container.go#L200
#   - Host : the dockerclient config element used to define the relationship
#        between the container and it's host.
#        https://github.com/fsouza/go-dockerclient/blob/master/container.go#L472
#
###

# This is an example service to show off some options
#
# @NOTE: the %INSTANCE token is created to match the running instance
example:
  Type: service
  Build: docker/ExampleImage # you will need ./coach/docker/ExampleImage/Dockerfile

  Instances: scaled 3 9  # start off with 3 running instances, allow up to 9

  Config:
    Image: myLocalImage:latest # Because there's a Build:, this will be built
    RestartPolicy: on-failure  # Not really needed

    Hostname: "%PROJECT_%INSTANCE"
    Domainname: "%CONTAINER_DOMAIN"

    Env:
      - "DNSDOCK_ALIAS=%PROJECT_%INSTANCE.%CONTAINER_DOMAIN"
      - "APP_KEY=%PERSONAL_APP_KEY"
      - "ENVIRONMENT=DEV"

    Entrypoint:
      - /app/.composer/vendor/bin/SomeAppBin
    CMD:
      - "--first-flag"
      - "--second-flag=VALUE"
      - "--third-flag=%CUSTOM_TOKEN"
    WorkingDir: /app/project/relative/path

    OpenStdin: true
    Tty: true

    ExposedPorts:
      3306/tcp: {}

  Host:
    Links:
      - OtherContainer:OtherContainer.local
    Binds:
      - "app/source:/app/source"
      - "app/assets:/app/assets"
      - "/tmp/absolute/path:/tmp/absolute/path"
    VolumesFrom:
      - OtherContainer

    PortBindings:
      80/tcp:
        - HostPort: 8080            # Port 8080 applies to all Host IPs

`)
  tasks.AddFile("app/README.md", `# Bare Project
## /.coach/

  Project coach configuration

## /app

  Project source-code and assets path

`)
  tasks.AddFile("app/source/README.md", `# Application source code root`)
  tasks.AddFile("app/assets/README.md", `# Non-Volatile assets root`)

  tasks.AddMessage("Created local project as a `starter` project")

  return true
}
