# coach/libs

The core libraries for the coach tool.

The libs define interfaces and standard implementations of client, nodes and instances.

## client

The client is responsible for connecting to a backend to manage images and containers

### factory

A client needs to implement a client factory to create clients per node, using a centralized
configuration in the client.yml file.

### fsouza-docker

The default (and currently the only) backend client option is the wrapper for the fsouza docker
library, which runs docker commands, using node and instance settings

## node

A node is an atomic configuration for an image and a set of containers for a single functional
service.
There are four different types of nodes

### build

A build node provides only an image, on which other nodes can depend.

### volume

A volume node provides a non-runnable container which other nodes can use for reliable file
space independent of platform.  It is a good idea to rely on volume nodes for source and asset
files so that an environment can be re-created in any format in any environment.

### service

A service node provides a running program, usually one that provides a tcp/ip socket.

### command

A runnable, disposable container setup that can be used to run a command as though it was a
local command.

## instances 

Instances are collection of instance struct, with particular behaviours.

### single

Singles Instances nodes provide nodes with a single container instance per node.

### fixed

Fixed instances nodes provide a fixed set of names instances/containers per node

### scaling

Scaling instances provide numerically names instances/containers per node, with identical
configurations which can be scaled up/down.

### temporary

Temporary instances provided disposable instances/containers for nodes, typically used for
command nodes.

## instance

An instance is a single object which correlates to a container.
