#!/bin/bash
. ./scripts/util.sh

set -e

optionalArgs=$1
package=$2
verbose=$3

if [[ "NO_CONSOLE_COLORING" == "true" ]]; then
    go test -v -p 1 -count=1 -timeout 1800s --tags development $optionalArg ./$package/... || exit 1
else
    go test -v -p 1 -count=1 -timeout 1800s --tags development $optionalArg ./$package/... \
        | awk '{gsub(/RUN/, "\033[34m&\033[0m"); gsub(/PASS|SUCCESS/, "\033[32m&\033[0m"); gsub(/FAILURE|FAIL|ERROR|error|Error|panic|panicked|SIGSEGV/, "\033[1;31m&\033[0m"); gsub(/WARN|WARNING|warn|warning|Warn|Warning/, "\033[1;33m&\033[0m"); gsub(/INFO|info/, "\033[36m&\033[0m"); gsub(/DEBUG|debug/, "\033[36m&\033[0m")} 1'
fi
