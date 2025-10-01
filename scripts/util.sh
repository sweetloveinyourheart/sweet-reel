#!/bin/bash

set -e
set -o pipefail

[ -f ./scripts/util_logs.sh ] && source ./scripts/util_logs.sh
[ -f ./scripts/util_docker.sh ] && source ./scripts/util_docker.sh
