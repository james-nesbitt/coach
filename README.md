# COACH

Coach is a docker-compose port to go, with an augmented yaml syntax.

The main purpose of coach is to provide a Docker based application
environment, which is composed of multiple containers, each of which
has it's own purpose.
Some of the containers exist to provide services, and can be scaled.
Some of the containers exist only to provide common builds, or shared
file/volume space.
Some of the containers are disposable command containers.

### Coming from docker-compose

Coach allows:
- defining some nodes as being of different types
- some node types can be used just for builds, or just for volumes
- some nodes run services in the background such as services
- some nodes can be run from the host as utility commands,

* coach has a stubbed out docker-compose.yml interpreter, that I just
  never got around to fleshing out.

Coach handles instances slightly differently:
- some nodes can be considered "scalable"
- some nodes can have only a single instance
- coach is aware of other node names for things like links and volumes (no 
  need to worry about container names when linking to other nodes)
- nodes can hardcode Links, VolumesFrom etc to particular instances

Coach also gives you
- a token system, to define a variable once, and use it across the configuration
- some tokens are kept per user, some are kept separate for secrecy
- some tokens are defined per container instance
- coach can init a project (kind of like git), and can do so using templates
- some coach templates are avaialble online
- you can build your own templates, and keep them wherever you want

## How to use Coach

### A quick breakdown
- Coach is a Go tool, and is usually installed using "go get"
- Coach is run as a command line tool, in any part of your project,
  which is demarked by a .coach folder (like .git.)
- A Coach project has a number of nodes, which are any number of
  containers that will run from the same image, with the same
  configuration
- Coach has a number of operations, which can be targeted to any
  of the project nodes

### Installing

#### You will need Go, so install go as per your environment

1. install go
2. include the go bin folder in your path, or use coach using 
  the full path.

#### You can retrieve the coach source:

    $/> go get github.com/james-nesbitt/coach

  * I will work on an installer when I get a chance, but currently you have to
    use the go method.

### Creating/Initializing a project
Starting a coach project means defining a certain folder as being the root of
your project.  Coach wants to see a .coach folder in that path. Coach has an
operation that it can use to "initialize" a project root for you, similar to 
how git has it's "init" operation.

If you know how to use coach, then start a bare project, and modify the 
new configurations in the .coach folder

If you are not sure how to configure coach, consider starting from a demo
template, which has working configurations for real world scenarios.

#### Starting a Bare empty project
Create a new "empty" local coach project in the current path

    $/> coach init
    
* This project will not have any runnable nodes, so the operations will
      have no targets to act on.

#### Starting a project with a bit more
Create a new coach project in the current path, with some example configurations

    $/> coach init starter
    
 * This project will have only a single "example" node built in, which will
   not run properly, as it doesn't use a real image.

#### Starting a project using one of the demo templates:

The Lamp Demo variations:

##### lamp
A standard, distributed service project, with different service containers for nginx,
php and mariadb

    $/> coach init demo lamp  

##### lamp singleton (1 container approach)
This demo runs a LAMP stack all in a single container, to keep the number of containers
to a minimum.  It is simple to use, and managed using supervisord.

    $/> coach init demo monolithic

##### lamp with multiple PHP versions running (PHP56, PHP7, HHVM)
This demo creates a project with a single nginx service, but multiple concurrently
running PHP services that can be accessed using their own URLs.

    $/> coach init demo multiplephp

#### The feature complete demos:
A kind of coach-feature demo that shows a number of configuration options:

    $/> coach init demo complete

A drupal8 from composer demo, with working composer, drush and console commands

    $/> coach init demo drupal8

### Configuring the created project
Once you have a .coach folder in place, you can configure the project by editing the
various files in the .coach folder.  The files created in the demos tend to have
plenty of documentation in their example files.  More information can be found
in the "coach help" and in the wiki.

#### .coach/conf.yml
This file contains various top level configurations, including a string map of tokens
that you can use for string substitution in the other files

#### .coach/secrets/secrets.yml
A token string map, that can be used for substution in other files, that is kept
separate from the conf.yml, so that it can container sensitive tokens, that you
may want to keep out of a revisioning repositoty.

#### .coach/nodes.yml
A string map of coach "nodes"

Each node maps together a node image, some container settings, and a host implementation
to allow the creation of one or more local containers, called "instances".
- In the case of build nodes, the node defines just a docker build
- volume nodes define single containers that contain volumes, but will never be started
- service nodes define scaleable containers that will run services
- command nodes define disposable command-line run containers

#### Other files
- .coach/help.yml : a custom project string map of help topics
- .coach/tools.yml : (read the wiki about tools)

### Daily Use
Typically used by a developer or a sysadmin, the coach system allows the creation of a project layout that is both friendly to a developer, but also easily used in production environments.
Most of the instructions here are targeted at the developer, but the same approach can be taken by a sysadmin.

Successful developer and systadmin approaches are definitely feasible with coach, but rely on a certain style of Docker usage, which makes it much more straight-forward to use, and simpler for deployment.

##### Local use
A developer typically uses docker locally, and wants to map their host source into a project.

The best approach here is to create a local volume node, and map/bind in host folders to the container. Other nodes then mount volumes from the volumes container.  Such a container could be recreated in production as a static container with source built in.

##### Remote docker use
In a remote docker scenario, the options to locally bind are not there, but the same volume container can be used, however with source built into it.  Then a remote-access service container could be used to give ssh/nfs/smbfs access to the container, and it can be used as a remote server.

Building source into a volume container uses a custom local image build, which will copy, or export (git) source code during the build.  Such an image would then need to be distributed.

To configure a project to use a remote docker client, you can either configure your docker to use a remote client (using ENV variables) or you can explicitly configure the docker client in
the project conf.yml

    #.coach/conf.yml
    Docker
        Host: "tcp://10.0.42.1"         # point to a remote docker server

#### Starting the first time

To get things started after initialization, you can use a single
operation, similar to using docker-composer:

    $/> coach up

  This will:

  1. build any docker images used
  2. pull any docker images not being built
  3. create any peristant containers
  4. start any service containers

  * Note that building and pulling images will overwrite any local
    versions, so consider avoiding those operations if you need to
    test out locally custom built or altered images

#### Starting manually

The manual equivalent of the "up" operation is the following sequence:

    $/> coach build
    $/> coach pull
    $/> coach create
    $/> coach start

#### Starting and stopping

After the initial start, images and containers should be created, and
the project can be started and stopped using only the two operations:

    $/> coach start

and

    $/> coach stop

The containers are persistant across these operations, so no data should
be lost.

Docker also provides a pause operation which can be used:

    $/> coach pause

and

    $/> coach unpause    

#### Rebuilding the coach project
Sometimes container break, or need to be tested from scratch, and sometimes images need updating.

#### Rebuilding images

You can rebuild on top of the existing image using the build operation

    $/> coach build

And you can update the remote images using pull

    $/> coach pull

#### Rebuilding containers

To rebuild containers, you must first stop and remove the previous ones:

    $/> coach stop
    $/> coach remove

Then you can re-create them, and start them again

    $/> coach create
    $/> coach start

### Advanced Use

#### Targeting operations at specific containers:

If you want an operation to target specific nodes, you can list them
before the operation, with an "@" in front of the name:

    $/> coach @www @fpm remove

You can also target all nodes of a certain type

    $/> coach %service create

If you are using a node that runs mulitple container instances, then you
can target specific instances

    $/> coach @www:3 stop

  * This usually requires that a node in your nodes.yml file has "Instances:"
    defined.

### Advanced operations
As docker includes many more operations, coach also tries to implement them

To get help on an operation
- help : get help information about coach, and coach operations

      $/> coach help
      $/> coach help {topic}
      $/> coach help {operation}

- info : give information about images and containers

       $> coach info

The full operations list:
- attach: attach to a running instance/container to see it's output
- build: build images for any nodes with a "Build:" declaration (requires a local build)
- clean: stop and remove container, remove any built images
- commit: commit a container to an image:revision
- create: create containers for a node
- pause: pause all processes in containers
- pull: pull images from nodes that use remote images
- remove: remove any non-running, created containers
- run: run a command in a new container for a node (not exec)
- scale: increate the number of running instances/containers of scaling nodes
- start: start any service containers
- stop: stop any running service containers
- unpause: unpause any paused containers

The non-standard operations are
- init: initialize a new coach project in the current path
- tool: run scripts or other external tools, but pass in coach tokens and container names

### Starting your own init templates
All of the core demos are run using a yaml syntax, that directs the init
operation to create files, copy files, checkout files and display messages.

For the core demos, those yaml files are pulled directly from the github 
coach project source, but coach can just as easily use any yaml file that
you create, be it local or retrievable on the internet.

The syntax for using your own local yaml file:

    $/> coach init yaml /path/to/my/init.yaml
    
or remote yaml

    $/> coach init yaml https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/drupal8/.coach/coachinit.yml

Take a look at the existing demos for the syntax:

- complete: https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/complete/.coach/coachinit.yml
- drupal8: https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/drupal8/.coach/coachinit.yml
- lamp: https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/lamp/.coach/coachinit.yml
- lamp_monolithic: https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/lamp_monolithic/.coach/coachinit.yml
- lamp_multiplephps: https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/lamp_multiplephps/.coach/coachinit.yml
  
- lamp_scaling: https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/lamp_scaling/.coach/coachinit.yml

### Coach tools
Sometimes using coach to run containers does not provide enough funcationlity for a project,
or sometimes it takes too much effort to use, whereas a simple script might suffice.  Coach
can integrate simple scripts, which allows you to pass coach tokens, and metatadata into a
script as environment variables; and it allows a script to be run anywhere in the coach 
project without worrying about relative paths.
