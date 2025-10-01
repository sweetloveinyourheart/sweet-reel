#!/bin/bash

# Console colors
red='\033[0;31m'
yellow='\033[0;33m'
green='\033[0;32m'
gray='\033[0;90m'
cyan='\033[0;36m'
nc='\033[0m' # No Color, use this to terminate color sequences

app_timestamp() {
    date +"%H:%M:%S"
    return 0
}

app_echo() {
    local timestamp=$(app_timestamp)
    local callLocation=$(caller)
    local callLocationFile=$(echo $callLocation | cut -d' ' -f2)
    local callLocationLine=$(echo $callLocation | cut -d' ' -f1)
    echo -e "${gray}[$timestamp][$callLocationFile:$callLocationLine]:${nc} $@"
    return 0
}

app_echo_red() {
    local timestamp=$(app_timestamp)
    local callLocation=$(caller)
    local callLocationFile=$(echo $callLocation | cut -d' ' -f2)
    local callLocationLine=$(echo $callLocation | cut -d' ' -f1)
    echo -e "${gray}[$timestamp][$callLocationFile:$callLocationLine]:${nc} ${red}$@${nc}"
    return 0
}

app_echo_yellow() {
    local timestamp=$(app_timestamp)
    local callLocation=$(caller)
    local callLocationFile=$(echo $callLocation | cut -d' ' -f2)
    local callLocationLine=$(echo $callLocation | cut -d' ' -f1)
    echo -e "${gray}[$timestamp][$callLocationFile:$callLocationLine]:${nc} ${yellow}$@${nc}"
    return 0
}

app_echo_green() {
    local timestamp=$(app_timestamp)
    local callLocation=$(caller)
    local callLocationFile=$(echo $callLocation | cut -d' ' -f2)
    local callLocationLine=$(echo $callLocation | cut -d' ' -f1)
    echo -e "${gray}[$timestamp][$callLocationFile:$callLocationLine]:${nc} ${green}$@${nc}"
    return 0
}

app_echo_blue() {
    local timestamp=$(app_timestamp)
    local callLocation=$(caller)
    local callLocationFile=$(echo $callLocation | cut -d' ' -f2)
    local callLocationLine=$(echo $callLocation | cut -d' ' -f1)
    echo -e "${gray}[$timestamp][$callLocationFile:$callLocationLine]:${nc} ${cyan}$@${nc}"
    return 0
}
