# DEMO lamp

## /app/www

The final web root to be used by most of the containers.

This gives a source build destination, as a subfolder to 
a bound path, which can be deleted, copied and moved as a
part of any source build process.

This separates writeable source into elements that an be
considered read-only, which are used in the nginx and 
fpm service.

## /app/source

This folder contains any source files and scripts which 
you may use to build your www path.
The build process may be run outside of the coach containers
(in which case there is no need to map it in) or it can
be mapped in, only to containers that build source.


## /app/assets

This mapped folder contains all writeable assets that need
to be integrated into your application. This folder allows
us to mount the www folder as read-only, in served cases,
for extra security


## /app/backup

An additional writeable path, into which the application
could archive and backup data.
