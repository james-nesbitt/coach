# COACH DEMO: Standard Lamp

This coach demo provides a standard LAMP stack:

- Nginx service using jamesnesbitt/wunder-nginx
- PHP-FPM service using jamesnesbitt/wunder-php56fpm
- DB service using jamesnesbitt/wunder-mariadb

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

    #/> coach init demo lampstack

  Now just do this:

    $/> coach up

  This should:

    - download some images
    - build any local images that were defined
    - create a source container which will map ./app/source to /app/www/active

    - start an nginx container, and fpm container and a db container

  You should then have

    - http access at host:8080, or directly on the container
    - the http should point to the source code in /app/source

# Day to day use

## starting and stopping

    $/> coach stop
    $/> coach start

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

# Alterations that you should make

## Custom DB

The nodes.yml file has a better DB image commented out which you could use to get
a DB container that has an already populated DB with specific access credentials.

## Switch to PHP7 or HHVM

If you want to test out the latest PHP7 or HHVM builds, try changing nodes.yml, the 
fpm node Image values.  There are instructions above the node.

Note that if you want to swap out fpm containers, you will need to stop and remove
both the fpm and nginx nodes, alter the nodes yml, and then start www and fpm again
(you do not need to stop the db container.)

Instructions:

1. Make changes to your nodes.yml (change the fpm image)

2. stop www and fpm

  $/> coach @fpm @www stop

3. remove www and fpm containers

  $/> coach @fpm @www remove

4. make sure that you have all needed images
 
  $/> coach pull
  $/> coach build

5. re-create www and fpm containers

  $/> coach @fpm @www recreate

6. restart www and fpm containers

  $/> coach @fpm @www start
