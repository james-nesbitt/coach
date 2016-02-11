# DEMO Drupal 8

## /app/drupal-project

Drupal 8 source as provided by the drupalproject

https://github.com/drupal-composer/drupal-project.git

This source is cloned during project init, so it is a
clean copy.  You will probably want to remote the git
remote, and make some changes

Note that many of the paths in the source are mapped
over in the containers, so that the assets paths can 
be mapped over the source code.  this allows source 
to be mounted as read-only, for security.

# /app/assets/

## /app/assets/files

The two folders inside this path will be mounted into
containers for the sites/default public and private 
files paths.

# /app/settings

Various cli app configurations

# /app/settings/console

Configurations for the Drupal console command container

# /app/settings/drush

Configurations and aliases for drush

# /app/settings/settings.local.php

An example for how to mount a local file over the source
code, this keeps the local settings file out of the source code

# /app/settings/services.local.yml

An example for how to mount a local services file over the source
code, this keeps the local services file out of the source code

# /app/backups

Mounted into the containers as an optional destination for
backup and archive dumping.
