# COACH DEMO: Drupal8 w/ composer

This coach demo provides a Drupal 8 site, using composer and 
the drupal-project composer repo.  The Drupal source code is
retrieved from git, so composer will need to be run before 
using Drupal, but instructions below should cover that.

The demo provides the following:

- standalone nginx and php7 fpm containers
- a standalone db server, with an empty db (mysql://app:app@db.app/app)
- read-only source for fpm (but read-write in dev containers)
- writeable assets, mapped directly into the sites default paths (no sym-links required)
- command containers for composer/console/drush

- a host mounted backup folder that can be used for drush dumps

The database is empty, but you can either include provisioning
in the custom DB build, or you can install into it.

There is a writeable settings.local.php in the app/drupal
folder, but you will have to play around to get settings.php
writeable (just duplicate what I did for the local settings)

All of the suggested images in this project are based on
the wunder-base concept, which is used to provide a commong
/app project folder, and similar users/OS, to keep access
privileges working better.
The builds are all CentOS7, EPEL and often REMI based builds.

# Getting Started:

  You already did this:

    #/> coach init demo drupal8

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

  As Drupal8 needs composer to run, you will need to use the composer command
  container to do the initial installs

    $/> coach @composer run update
    $/> coach @composer run install

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

  @NOTE that you will lose your DB contents if you don't export them

  1. stop any running containers

    $/> coach stop

  2. remove and created containers

    $/> coach remove

  3. recreate and containers

    $/> coach create

  4. start the new containers

    $/> coach start

# Alterations that you should make

## Installable

If you can get drush to install, it uses the source as read-write
so it could avoid any problems.

You could add a writeable settings.php on top of the read-only
source, to allow a Drupal install.  Really this should not be 
necessary, as there are plenty of patches out there to allow 
Drupal to be installed if the settings.php give a DB connection
already.

You can add a mount to the source container:

  Config:
    Volumes:
      "/app/project/web/sites/default/settings.php": {}

But you will have to copy the default.settings.php into that new
file

## Custom DB

You could alter the DB build to include an sql dump
