# To generate this file
#   1. if not "chmod +x gen.sh", do it
#   2. "./gen.sh" at /proto path

#!/bin/bash

# Ensure the script is executable with: chmod +x gen.sh

# Set default paths using `which` to find installed tools
PROTOC_GEN_GO=$(which protoc-gen-go)
PROTOC_GEN_GO_GRPC=$(which protoc-gen-go-grpc)

# Check if the plugins are installed
if [ -z "$PROTOC_GEN_GO" ]; then
  echo "Error: protoc-gen-go not found. Please install it and ensure it's in your PATH."
  exit 1
fi

if [ -z "$PROTOC_GEN_GO_GRPC" ]; then
  echo "Error: protoc-gen-go-grpc not found. Please install it and ensure it's in your PATH."
  exit 1
fi

# Generate Go code for server
protoc --go_out=../server --go-grpc_out=../server event.proto

# Generate Go code for client
protoc --go_out=../client --go-grpc_out=../client event.proto

echo "Protobuf generation complete."

