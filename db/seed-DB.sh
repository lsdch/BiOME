#! /bin/sh

echo "Database will be wiped before inserting data seeds."

edgedb database wipe

echo "Migrating database schema... (first migration might take a moment)"
edgedb migrate --schema-dir schema

echo "Importing initial data"
cd ./seeds
go run main.go
