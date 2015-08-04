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
    mysql -uroot -e "GRANT ALL ON project.* to project@'172.17.%' IDENTIFIED BY 'project'" && \
    mysql -uroot -e "FLUSH PRIVILEGES"

### /ProjectDB -------------------------------------------------------------------
