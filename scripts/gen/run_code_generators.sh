#!/bin/bash
. ./scripts/util.sh

set -e

function resetFiles() {
    app-echo "Resetting files matching $1"
    local filePattern=$1
    local files=$(find proto -name "$filePattern" -type f)
    for file in $files; do
        git checkout --ours $file &> /dev/null || ( git checkout --theirs $file &> /dev/null || : ) # :)
    done
}

resetFiles "*.pb.go"
resetFiles "*.connect.go"

goImportsCmd="go run golang.org/x/tools/cmd/goimports --local "github.com/sweetloveinyourheart/sweet-reel" -w ./"
goGenerateCmd="go generate --tags generate ./..."

app-echo "Running goimports..."
$goImportsCmd

app-echo "Running go generate..."
$goGenerateCmd || (app-echo "go generate failed, retrying after goimports..." && $goImportsCmd && $goGenerateCmd)