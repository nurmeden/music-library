#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <migration_name>"
    exit 1
fi

MIGRATION_NAME=$1
TIMESTAMP=$(date +"%Y%m%d%H%M%S")
FILENAME="${TIMESTAMP}_${MIGRATION_NAME}.sql"

cat <<EOL > "migrations/$FILENAME"
EOL

echo "Migration file created: migrations/$FILENAME"
