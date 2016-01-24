# COACH

A container management tool, written in go, that can wrap the fsouza docker library.

The primary implementation is a command line tool (the cli) which can be used to create,
start and stop various containers in groups. Running coach typically involves running 
command line commands to prepare an environmment, and to start up containers.

Originally, this started as a go port of the docker-compose application, this application
grew to be a robust container management tool, with an astracted docker backend, and a 
comprehensive modular approach.
Coach primarily uses docker as a backend, and can be configured to use a remote docker 
client if needed.

The tool is easy to use, and has integrated help, but requires some learning to understand
configuration files.  To make things a bit easier, there is an initialization tool that
can also install a number of demo configurations.

Overall, the tool is stable, and realiable, although it could still do a bit more checking
on operations before reporting warnings.

There is a wiki in the git, which can be used to get comprehensive information.

## Architecture

Coach was writtent primarily to work as a cli tool, in a project context.  The coach 
tool recognizes a project context by looking for a .coach folder, in the same way that
the git cli looks for a .git folder.
The .coach folder marks the root of the project and also contains a number of files which
can be used to configure the project.  Coach also looks for files in ~/.coach, to allow 
for configuration elements to be user centric, but it still needs a project .coach folder 
to mark the root of a project.

The project root is important as it is the root path used for things like local file
volume mounts.  The .coach folder is important because it can contain configuration files
and docker builds.
Typically in a project, any coach specific files can be kept in the .coach folder, to 
keep coach specific elements out of your project code.

The coach code itself is written as a linear (single-threaded) set of libraries, with 
a cli wrapper.
The primary coach concepts are:

- client : a backend library that is responsible for creating, removal and starting
   as well as stopping of docker images and containers
- client-factory : a configuration that is meant to create clients

- project : a set of configurations and paths for a coach project, including some string
  tokens used for string replacements
- tokens : a string map that can be used for string replacement.
- paths : a path map

- nodes : a map of nodes for a project
- node : a single configuration which is used to create a single image, and a set of 
  containers to provide a docker build, a service, some volume binds, or a command
- instances : a set of instances for a node, with a particular behaviour.
- instance : a particular container for a node

- operation : a cli task, that runs node client functions for a project.

## Installing

Before you can use coach you will need to install go in your working environment. 
Installing go is different in every environment, but usually just involved installing
a go binary, and then configuring a go root.
If you don't configure a go root, then your machine will not be able to find 

Once go is installed, you can install the coach tool using the following go command:

    $/> go get github.com/james-nesbitt/coach

This will download a number of packages, including the following:

- https://github.com/james-nesbitt/coach : the coach tool
- https://github.com/fsouza/go-dockerclient : a go library for the docker remote API
- https://github.com/go-yaml/yaml : a go library for yaml parsing
- github.com/twmb/algoimpl/go/graph : a go library for graph sorting (used to order nodes)

## Quick start

To use coach you will need a project, with a project root and a .coach folder in that root.
The .coach folder should contain a conf.yml file, and a nodes.yml file.  If you want to 
use a non-default backend client, then a client.yml file can be used.

### Init

The easiest way to get a working .coach installation is to use the init command.  You can
use the init directly through the coach cli, which will install an empty configuration set
that doesn't really do anything:

    $/> coach init

A slightly more usefull example installation can be installed using the starter init:

    $/> coach init starter

After using the init command, alter the .coach yaml files to properly configure your project.

### Demos

The init command can also install some prepared demonstration examples as coach projects.
Using a demo will install an unnamed project into the current path.  Use the coach init
command and alter the configurations to meet your needs.

Some more insteresting demo configurations can be tested using the following:

- A single LAMP container

    $/> coach init demo lamp_singleton

  This demo provides a single docker container that runs a set of LAMP services.

- A single LAMP split container set

    $/> coach init demo lamp

  This demo uses a number of containers to give a LAMP service, with each container
  running a single service.  This is a best-practices docker approach.

- A single LAMP container with various php server options

    $/> coach init demo lamp_multiplephps

- A scaling www/fpm server set

    $/> coach init demo lamp_scaling

### Alternate initializations

The demo configurations run using yml configuration files that are kept in the coach
repo.  The yml format used for the demos can be used by anyone to create default installations
that anyone with access to the files can use.  This provides easily distributable installations
for coach usage.

Alternatively, there is no reason why the coach configuration for a project cannot be 
included in a project repository if desired.  It is well encapsulated, and also provides 
resources for allowing user specific overrides of project configurations, so that users can use
their own configuration values. and so that sensitive information does not need to be distributed.
