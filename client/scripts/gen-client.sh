#! /env/bin/bash

if [ \"$(md5sum openapi.json | cut -f1 -d' ')\" != \"$(cat openapi.json.md5)\" ]; then
  echo $(md5sum openapi.json | cut -f1 -d' ') >openapi.json.md5
  openapi-ts -i openapi.json -o src/api
else
  echo "✅ No changes in API spec"
fi
