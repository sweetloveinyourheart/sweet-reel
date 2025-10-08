#!/bin/bash
. ./scripts/util.sh

set -e

function resetFiles() {
    app_echo "Resetting files matching $1"
    local filePattern=$1
    local files=$(find proto -name "$filePattern" -type f)
    for file in $files; do
        git checkout --ours $file &> /dev/null || ( git checkout --theirs $file &> /dev/null || : ) # :)
    done
}

resetFiles "*.pb.go"
resetFiles "*.connect.go"

goGenerateCmd="go generate --tags generate ./..."
goImportsCmd="go run golang.org/x/tools/cmd/goimports@latest --local "github.com/sweetloveinyourheart/sweet-reel" -w ./"

app_echo "Running goimports..."
$goImportsCmd

app_echo "Running go generate..."
$goGenerateCmd || (app_echo "go generate failed, retrying after goimports..." && $goImportsCmd && $goGenerateCmd)