#!/usr/bin/env bash

set -eo pipefail

mkdir -p ./tmp-swagger-gen
cd proto

proto_dirs=$(find $(go env GOPATH)/src/github.com/cosmos/cosmos-sdk/proto ./network -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
  if [[ ! -z "$query_file" ]]; then
    buf generate --template buf.gen.swagger.yaml $query_file
  fi
done

cd ..
go-swagger-merger -o ./docs/static/openapi.yml $(find ./tmp-swagger-gen -name "*.swagger.json") ./docs/static/title.yaml

# clean swagger files
rm -rf ./tmp-swagger-gen