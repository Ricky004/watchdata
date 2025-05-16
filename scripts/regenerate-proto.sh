#!/bin/bash

set -e

PROTO_DIR="./proto/telemetry"
OUT_DIR="./"

echo "Regenerating Go protobuf files in $PROTO_DIR ..."

for file in $(find "$PROTO_DIR" -maxdepth 1 -name "*.proto"); do
  echo "Processing $file"
  protoc --go_out="$OUT_DIR" --go_opt=paths=source_relative "$file"
done

echo "Done."
