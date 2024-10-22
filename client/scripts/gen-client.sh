#! /env/bin/bash

force=false
while [ $# -gt 0 ]; do
  case $1 in
  -f | --force)
    force=true
    ;;
  esac
  shift
done
echo "current md5:    $(md5sum openapi.json | cut -f1 -d' ')"
echo "reference md5:  $(cat openapi.json.md5)"
if [[ \"$(md5sum openapi.json | cut -f1 -d' ')\" != \"$(cat openapi.json.md5)\" ]] || [[ $force == "true" ]]; then
  echo $(md5sum openapi.json | cut -f1 -d' ') >openapi.json.md5
  openapi-ts -i openapi.json -o src/api
else
  echo "✅ No changes in API spec"
fi
