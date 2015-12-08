# DEMO : Complete

The /app folder is a central place into which to put all of your
Application specific source code, configurations and assets, in
such a manner that this folder can be a distribution target for 
an application (keep this in git.)
Think of this path as being the root of your actualy project,
wheras the parent could be considered the root of the coach
implementation.  /app should be something that could work no 
matter which tool was being used to manage the project.

If your project is generally a .coach folder, then you can also 
keep that path in the same repository, but it is not compulsory.

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

Backup can be a target for manuall archivng, but it should be
considered a part of any sysadmin tools, and should be used
for scheduled and other automated backups.

## /app/settings

A place to put developer, or sysadmin configurations that may
get mapped into command containers, or used in specific
situations.
