#!/bin/bash
. ./scripts/util.sh

# This script is used to run all of the tests
# It is intended to be run from the Makefile

generateCov=0
verbose=0

if [[ "$1" == "verbose" || "$2" == "verbose" ]]; then
    verbose=1
fi

if [[ "$1" == "cov" || "$2" == "cov" ]]; then
    generateCov=1
fi

function runTest() {
    if [ $verbose -eq 1 ]; then
        $1 && ( app_echo_green "$1 - PASS" && exit 0 ) || ( app_echo_red "$1 - FAIL" && exit 1 )
    else
        export NO_CONSOLE_COLORING=true
        $1 &> tests/logs/log-$2 && ( app_echo_green "$1 - PASS" && exit 0 ) || ( app_echo_red "$1 - FAIL" && exit 1 )
    fi
}

mkdir -p tests/logs

if [ $generateCov -eq 1 ]; then
    if [ $verbose -eq 0 ]; then rm -rf tests/logs/log-cov-* && rm -rf tests/logs/log-lint; fi
    makeCommands=$(sed -n -e '/^cov-/p' Makefile | awk -F ":" '{print $1}')
else 
    if [ $verbose -eq 0 ]; then rm -rf tests/logs/log-ut-* && rm -rf tests/logs/lint; fi
    makeCommands=$(sed -n -e '/^ut-/p' Makefile | awk -F ":" '{print $1}')
fi

exitCode=0
runTest "make lint" lint || exitCode=1
for cmd in $makeCommands; do
    runTest "make $cmd" $cmd || exitCode=1
done

if [ $exitCode -eq 0 ]; then
    app_echo_green "All tests passed."
else
    app_echo_red "Some tests failed.  Run the individual make command(s) or review the logs for more details."
fi

exit $exitCode
