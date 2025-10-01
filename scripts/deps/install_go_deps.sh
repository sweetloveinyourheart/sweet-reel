#!/bin/bash
. ./scripts/util.sh

# This script installs go dependencies

# It is intended to be run from the Makefile

set -e

go mod download || go mod download || exit 1 # Try twice, in case of flakiness

go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest

exit 0
