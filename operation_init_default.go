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
source:
  Type: volume

  Config:
    Image: "jamesnesbitt/wunder-base"

  Host:
    Binds:
      - app/www:/app/www          # host based webroots folder (needs /active subroot for nginx conf)

# Database service
db:
  Type: service
  Build: docker/db               # DB has a docker build so that we can create databases and set custom passwords.

  Config:
    RestartPolicy: on-failure

    ExposedPorts:
      3306/tcp: {}

  Host:
    VolumesFrom:
      - files                              # I am not sure if this is needed

# FPM service
fpm:
  Type: service

  Config:
    Image: jamesnesbitt/wunder-php5fpm     # The FPM works, and should have blackfire working
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
www:
  Type: service

  Config:
    Image: jamesnesbitt/wunder-nginx
    RestartPolicy: on-failure

    Hostname: "%PROJECT_%INSTANCE"                  # Token : project name (can be set in conf.yml)
    Domainname: "%DOMAIN"                           # Token : can be set in conf.yml:Tokens
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

## Assets

The assets folder is meant to be a non-versioned folder that contains elements needed to
run the application, but which should not be a part of the project source code.  This includes
file assets, and cache elements and temprorary elements.
The real goal of this folder is to separate filespace into Read-Only and Writeable, with the
Assets folder being the writeable.

## Backups

`,
		"app/assets/README.md":  `

`,
		"app/backup/README.md":  `
The assets folder is meant to be a non-versioned folder that contains elements needed to
run the application, but which should not be a part of the project source code.  This includes
file assets, and cache elements and temprorary elements.
The real goal of this folder is to separate filespace into Read-Only and Writeable, with the
Assets folder being the writeable.
`,
		"app/www/active/index.php":  `<?php phpinfo();`,
	".coach/docker/db/Dockerfile": `
	FROM        jamesnesbitt/wunder-mariadb
	MAINTAINER  james.nesbitt@wunderkraut.com

	### ProjectDB --------------------------------------------------------------------

	# Create our project DB
	RUN (/usr/bin/mysqld_safe &) && sleep 5 && \
			mysql -uroot -e "UPDATE mysql.user SET Password=PASSWORD('RESETME') WHERE User='root'" && \
			mysql -uroot -e "DELETE FROM mysql.user WHERE User=''" && \
			mysql -uroot -e "DROP DATABASE test" && \
			mysql -uroot -e "UPDATE mysql.user SET Password=PASSWORD('RESETME') WHERE User='root'" && \
			mysql -uroot -e "CREATE DATABASE project" && \
			mysql -uroot -e "GRANT ALL ON project.* to project@'10.0.%' IDENTIFIED BY 'project'" && \
			mysql -uroot -e "FLUSH PRIVILEGES"

	### /ProjectDB -------------------------------------------------------------------
`,

	}
}
