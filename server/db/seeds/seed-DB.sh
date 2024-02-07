#! /bin/sh

echo "Database will be wiped before inserting data seeds."

edgedb database wipe

echo "Migrating database schema..."
edgedb migrate --schema-dir ../schema

go run ./seeds.go
