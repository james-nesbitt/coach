# COACH DEMOS

Each subfolder contains the base init template for a coach project
that can be used to demonstrate a set of functionality.

Each project is written to produce an immediately buildable and 
testeable set of containers, that should work, although in most
cases, only simple source code is used (such as the phpinfo() function)

## Architecture

Each demo is given it's own folder, to demonstrate the files.

### Using a demo

The simplest way to use a demo is to use the coach init method.

    $/> coach init demo {demo}

All of the builds include a remoteYaml build file which could be
directly used if desired:

    $/> coach init remoteyaml https://raw.githubusercontent.com/james-nesbitt/coach/master/templates/demo/lamp_multiplephps/.coach/coachinit.yml

## Demos

The following demos are available

### lamp

A simple LAMP stack, with each service in it's own container, and source mapped
into a separate container.

### lamp_multiplephps

This LAMP demo runs a single containers for DB and nginx but runs three PHP
servers in parralel, to allow testing of PHP56, PHP7 and HHVM on a single 
source code base. 

### lamp_monolithic

This LAMP demo uses a single LAMP stack container, with all services managed by
supervisord.

### lamp_scaling

This LAMP demo implements using a set of matched nginx-phpfpm servers to allow
testing of applications that need load-balance front ends.  The project by default
runs three pairs of nginx-php7 services, connected to a single db and source set.

## Advanced Demos

These demo options show off more complete coach features, but are less 
functionally tested, and may refer to application functionality that is not
included, such as "building your application from source"

### Complete

An attempt to get complete documentation and implementation of coach features
into a LAMP set.  This demo is extensive, but not deeply tested.

Features demonstrated:

- extensive configuration of coach
- secret/sensitive tokens
- custom shell script tool
- custom help topics
- custom Docker build
- volume nodes (including read only sharing)
- build base node (db)
- singular service nodes
- command nodes
- source separation
