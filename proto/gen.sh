# To generate this file
#   1. if not "chmod +x gen.sh", do it
#   2. "./gen.sh" at /proto path

protoc --plugin=protoc-gen-go=/Users/awats/go/bin/protoc-gen-go --go_out=../server event.proto
protoc --plugin=protoc-gen-go-grpc=/Users/awats/go/bin/protoc-gen-go-grpc --go-grpc_out=../server event.proto

protoc --plugin=protoc-gen-go=/Users/awats/go/bin/protoc-gen-go --go_out=../client event.proto
protoc --plugin=protoc-gen-go-grpc=/Users/awats/go/bin/protoc-gen-go-grpc --go-grpc_out=../client event.proto
