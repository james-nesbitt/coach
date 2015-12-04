# COACH DEMO: Monolithic LAMP

This coach demo provides a monolithic LAMP stack:

- server service using jamesnesbitt/wunder-lampstackplus

To incorporate source, the source code from /app/source
is mapped into a "source" container, which is used in
the nginx and fpm containers.

# Getting Started:

  You already did this:

    #/> coach init demo lamp_monolithic

  Now just do this:

    $/> coach up

  This should:

    - download an image

    - start nginx, fpm containers and db in a single container
      using supervisord

# Day to day use

## starting and stopping

    $/> coach stop
    $/> coach start

## rebuild the container

  1. stop the running container

    $/> coach stop

  2. remove container

    $/> coach remove

  3. recreate container

    $/> coach create

  4. start the new container

    $/> coach start
