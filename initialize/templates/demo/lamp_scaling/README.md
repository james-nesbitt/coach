# COACH DEMO: Scaling LAMP

This coach demo provides a scaling LAMP stack:

- a single DB service using jamesnesbitt/wunder-mariadb

- a scaling set of nginx/fpm pairs, starting at 3
  but scaling up to 9 instances.

  @these come in pairs to allow scaling http access, but
    DNS should work if you do some local magic.

To incorporate source, the source code from /app/source
is mapped into a "source" container, which is used in
the nginx and fpm containers.

The Database is un-initialized, but an optional init
build for a db has been included, and you can switch to
it by removing the db node, and uncommenting the custom
db node.

All of the suggested images in this project are based on
the wunder-base concept, which is used to provide a commong
/app project folder, and similar users/OS, to keep access
privileges working better.
The builds are all CentOS7, EPEL and often REMI based builds.

# Getting Started:

  You already did this:

    #/> coach init demo lamp_scaling

  Now just do this:

    $/> coach up

  This should:

    - download some images
    - build any local images that were defined
    - create a source container which will map ./app/source to /app/www/active

    - start some nginx containers, and  some fpm containers and a single db container

  You should then have

    - http access to all nginx services by directly pointing to the container

# Day to day use

## starting and stopping

    $/> coach stop
    $/> coach start

## starting and stopping more instances

  Use the scale operation to start additional instances:

    $/> coach scale up

  And the same operation to scale down

    $/> coach scale down

## starting a specific container

  - start a specific node

    $/> coach @www start

  - start all service nodes

    $/> coach %service start

## rebuild any containers

  1. stop any running containers

    $/> coach stop

  2. remove and created containers

    $/> coach remove

  3. recreate and containers

    $/> coach create

  4. start the new containers

    $/> coach start
