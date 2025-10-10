package grpc_videomanagement

// Run code generation with: go generate ./...
// gRPC documentation: https://grpc.io/
// To generate proto:
// - remove old *.pb.* files
// - add path to proto generation tools to PATH
// - tell buf: https://buf.build/ to generate proto in current directory with buf.gen.yaml instructions
//go:generate bash -e -o pipefail -c "rm -f *.pb.*; d=./../$(git rev-parse --show-prefix); cd $(git rev-parse --show-toplevel)/proto; export PATH=$(git rev-parse --show-toplevel)/.gobincache:$DOLLAR{PATH}; buf generate --template $DOLLAR{d}buf.gen.yaml $DOLLAR{d}"
