FROM        jamesnesbitt/wunder-mariadb
MAINTAINER  james.nesbitt@wunderkraut.com

### ProjectDB --------------------------------------------------------------------


# Create our project DB
#
# - set the root user password to something safe (considered a good idea)
# - create an "app" database
# - give access to the new database to an "app" user using the password "app"
# - app user can only access from 172.* IPs !!!! (inside the docker subnet)
# - flush access
#
RUN (/usr/bin/mysqld_safe &) && sleep 5 && \
    mysql -uroot -e "UPDATE mysql.user SET Password=PASSWORD('H4x0r') WHERE User='root'" && \
    mysql -uroot -e "CREATE DATABASE app" && \
    mysql -uroot -e "GRANT ALL ON app.* to app@'172.%' IDENTIFIED BY 'app'" && \
    mysql -uroot -e "FLUSH PRIVILEGES"

### /ProjectDB -------------------------------------------------------------------
