#! /bin/sh

echo "Database will be wiped before inserting data seeds."

edgedb database wipe

echo "Migrating database schema..."
edgedb migrate --schema-dir schema

echo "Importing initial data"
cd ./seeds
go run
