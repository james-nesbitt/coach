# COACH DEMO: Lamp with multiple PHP backends

This coach demo provides an ehnanced LAMP stack:

- DB service using jamesnesbitt/wunder-mariadb

- 3 fpm services, that connect a standard to a different 
  php version.  All php versions point to the same DB 
  container, and use a separate nginx host.
    - php56fpm : from centos REMI
    - php7fpm : from REMI-7
    - hhvm : from EPEL I think

- 1 master nginx server that has three server{} settings
  to point to each of the PHP FPMs using wildcarded URLS
    - php56.* -> the php56 fpm server
    - php7.* -> the php7 fpm server
    - hhvm.* -> the hhvm server

@NOTE you may have to do some DNS/hosts magic to get these
DNS entries to resolve

This should produce 3 parrallel nginx services that 
all point to the same source, which can be used for 
parrallel testing.

To incorporate source, the source code from /app/source
is mapped into a "source" container, which is used in
the nginx and fpm containers.

The Database is un-initialized, but it is recommended that
you include a local DB build, in which you create a DB, 
and add access parameters as needed.

All of the suggested images in this project are based on
the wunder-base concept, which is used to provide a commong
/app project folder, and similar users/OS, to keep access
privileges working better.
The builds are all CentOS7, EPEL and often REMI based builds.

# Getting Started:

  You already did this:

    #/> coach init demo lamp_multiplephp

  Now just do this:

    $/> coach up

  This should:

    - download some images
    - build any local images that were defined
    - create a source container which will map ./app/source to /app/www/active

    - start nginx containers, and fpm containers and a db container

  You can access the different http servers directly, or set
  different port numbers for each

# Day to day use

## starting and stopping

    $/> coach stop
    $/> coach start

## coach targets

  You can target coach at only certain nodes using targets

    $/> coach @www stop

    $/> coach @php56 @php7 @hhvm start

## rebuild any containers

  1. stop any running containers

    $/> coach stop

  2. remove and created containers

    $/> coach remove

  3. recreate and containers

    $/> coach create

  4. start the new containers

    $/> coach start
