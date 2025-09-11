#!/bin/bash

set -e
set -o pipefail

# Get migration name from argument
if [ -z "$1" ]; then
  echo "Usage: $0 <migration_name>"
  exit 1
fi
migration_name=$1

# Create a new migration file with a unix timestamp
timestamp=$(date +%s%N | cut -b1-13)

migration_file="internal/sql/migration/${timestamp}_${migration_name}.sql"
touch "$migration_file"

echo "Created migration file: $migration_file"
