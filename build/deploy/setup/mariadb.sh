#!/usr/bin/env bash

# pull container and install libs
docker pull tutum/mariadb
sudo apt-get -y install mariadb-client-5.5 libmariadbd-dev

# create a default MariaDB instance dba=admin@mariadb
docker create --name mariadb -p 3306:3306 -e MARIADB_PASS="mariadb" tutum/mariadb

if [ ! -L "/usr/local/bin/maria-db" ]; then
	sudo ln -s "$SETUP_HOME/mindcastio/deploy/lib/maria-db.rb" /usr/local/bin/maria-db
fi
