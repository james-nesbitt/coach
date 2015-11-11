package main

func (operation *Operation_Init) Init_Default_Run(flags []string) (bool, map[string]string) {
	return true, map[string]string{

		".coach/conf.yml":  `Project: coach
#Author: Used for docker commits

Paths:  # a map of paths and path overrides (see conf.go)
  test: "/test/path"

Tokens: # a map of string tokens, used for token replacement in the nodes.yml
  CONTAINER_DOMAIN: "docker"

Settings:
  UseEnvVariablesAsTokens: "yes"   # include all of the user's ENV variables as possible tokens

Docker:  # Override Docker configuration
#  Host: "tcp://10.0.42.1"         # point to a remote docker server
`,

		".coach/nodes.yml":  `
# Files volume container
#
# - A volume container used to hold files assets for an application
# - Also has a place to put backups (separate from file assets)
# - Volatile, and not likely to handle exports well
#
# * Should be used as a ReadWrite container for Links/VolumesFrom
#
files:
  Type: volume

  Config:
    Image: "jamesnesbitt/wunder-base"
    Volumes:
      "/app/tmp": {}               # /app/tmp is a volatile container folder
  Host:
    Binds:
      - app/assets:/app/assets     # host based assets folder
      - app/backup:/app/backup     # host based archive folder

# Source volume container
#
# - A volume container to hold application source
#
# * Can be used as a ReadOnly container for Links/VolumesFrom
#
source:
  Type: volume

  Config:
    Image: "jamesnesbitt/wunder-base"

  Host:
    Binds:
      - app/www:/app/www          # host based webroots folder (needs /active subroot for nginx conf)

# Database service
#
# - Standalone DB server
#
# * ExposedPorts is likely not necessary, it just gives a host port for the server
#
db:
  Type: service
  Build: docker/db               # DB has a docker build so that we can create databases and set custom passwords.

  Config:
    RestartPolicy: on-failure

    ExposedPorts:
      3306/tcp: {}

# FPM service
#
# - Standalone php-fpm service
#
# * iIt needs source, and assets 
# * ExposedPorts is likely not necessary, it just gives a host port for the server
#
# ! Alternate image: jamesnesbitt/wunder-php7fpm
# ! Alternate image: jamesnesbitt/wunder-hhvm
#
fpm:
  Type: service

  Config:
    Image: jamesnesbitt/wunder-php56fpm     # The FPM works, and should have blackfire working
    RestartPolicy: on-failure

    ExposedPorts:
      9000/tcp:

  Host:
    Links:
      - db:database.app
    VolumesFrom:
      - files
      - source

# WWW service
#
# - Standalone nginx service
# - Sets Hostname and DomainName using coach tokens
# - Sets DnsDock Alias ENV var using coach tokens
# - Tries to bind to the Host 8000 port
#
# * It needs source, and assets 
# * ExposedPorts are likely not necessary, they just gives a host port for the server
# * If you make this scaled, you will get an error trying to reused the 8080 port
#
www:
  Type: service

  Config:
    Image: jamesnesbitt/wunder-nginx
    RestartPolicy: on-failure

    Hostname: "%PROJECT_%INSTANCE"                  # Token : project name (can be set in conf.yml)
    Domainname: "%DOMAIN"                           # Token : environment domain (can be set in conf.yml)
    Env:
      - "DNSDOCK_ALIAS=%PROJECT.%CONTAINER_DOMAIN"  # If you are using DNSDOCK, this will create a DNS Entry.

    ExposedPorts:
      80/tcp: {}
      443/tcp: {}

  Host:
    Binds:
      - app/www:/app/www
    Links:
      - fpm:fpm.app
    VolumesFrom:
      - files
      - source
    PortBindings:
      80/tcp:
        - HostPort: 8080            # Port 8080 applies to all Host IPs

`,

		".coach/secrets/secrets.example.yml":  `# SECRET TOKENS THAT CAN BE KEPT OUT OF GIT
# Move this file to secrets.yml for it to be used
SECRET: VALUE
SECONDSECRET: "OTHER VALUE"
THIRDSECRET: "%PROJECT.SOMETHING"
`,

		"app/README.md":  `
# The coach app folder

The purpose of this folder is to keep all of the project elements that in a single location,
in a manner that maps easily into the wunder-base approach.  This approach is meant to keep
all source code mapped well into the containers, for both path and ownership.

## WWW

The www folder is meant to house the htdocs parts of the apllication.  This gives a render/make
target for application source code, into which it can be linked or copied.  This allows separation
of project custom code, from community code for frameworks and libraries.

* The typical nginx configuration expects an application Web Root at www/active

## Assets

The assets folder is meant to be a non-versioned folder that contains elements needed to
run the application, but which should not be a part of the project source code.  This includes
file assets, and cache elements and temprorary elements.
The real goal of this folder is to separate filespace into Read-Only and Writeable, with the
Assets folder being the writeable.

## Backups

The backups folder is meant to be a non-versioned fodler that contains backups dumps for
the application, which are kept separate from assets to that they can be managed separately.
`,
		"app/assets/README.md":  `
The assets folder is meant to be a non-versioned folder that contains elements needed to
run the application, but which should not be a part of the project source code.  This includes
file assets, and cache elements and temprorary elements.
The real goal of this folder is to separate filespace into Read-Only and Writeable, with the
Assets folder being the writeable.
`,
		"app/backup/README.md":  `
The backups folder is meant to be a non-versioned fodler that contains backups dumps for
the application, which are kept separate from assets to that they can be managed separately. 
`,
		"app/www/active/index.php":  `<?php 
/**
 * This is your web root, where you application root should sit.
 *
 * This is kept as a sub-path of the www folder, so that the www folder can be mapped into
 * the container, with a sub-folder that can be the target of a build process which may replacement
 * or sym-link it.
 * If you are using a build process to generate your web-root, have it build the active folder,
 * or symlink the latest build to ./active.
 **/
phpinfo();`,
	".coach/docker/db/Dockerfile": `#####
# Create a custom DB build for my project
#
# - this gives a custom DB for my project
# - my custom DB will respond to credentials that I define here
# - this DB image can be used as a "docker commit" target to create DB snapshots
#
# * This image build was created by coach init

FROM        jamesnesbitt/wunder-mariadb
MAINTAINER  james.nesbitt@wunderkraut.com

### ProjectDB --------------------------------------------------------------------

# Create our project DB
#
# - create a new database "app"
# - grant privileges to that database to user app@172.*, using the password "app"
# - flush privileges
#
# * this expects Docker to run it's bridge/subnet using 172.* (and limits access to that subnet)
# 
RUN (/usr/bin/mysqld_safe &) && sleep 5 && \
		mysql -uroot -e "GRANT ALL ON app.* to app@'172.%' IDENTIFIED BY 'app'" && \
		mysql -uroot -e "FLUSH PRIVILEGES"

### /ProjectDB -------------------------------------------------------------------
`,

	}
}
