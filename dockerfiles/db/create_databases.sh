#!/bin/bash

set -e
set -u

function create_user_and_database() {
    local database=$(echo $1 | tr ':' ' ' | awk  '{print $1}')
	local user=$(echo $1 | tr ':' ' ' | awk  '{print $2}')
	echo "  Creating user and database '$user' : '$database'"
	createuser -U $POSTGRES_USER $user --superuser;
	createdb -U $POSTGRES_USER $database;   
}

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
	echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
	for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ';' ' '); do
		create_user_and_database $db
	done
	echo "Multiple databases created"
fi
