#!/bin/bash

app_docker_cleanup() {
    pkill compose &> /dev/null || : # compose logs processes are sometimes left running, kill those too
    docker network rm -f fs_test_net &> /dev/null || :
    docker network prune -f &> /dev/null || :
    docker container prune -f &> /dev/null || :
    docker system prune -f &> /dev/null || :
    return 0
}

app_compose_up() {
    local composeFile=$1

    if [ -z "$composeFile" ]; then
        app_echo_red "call requires a composeFile"
        return 1
    fi

    composeCommand="docker compose -f $composeFile up --wait-timeout 300"
    app_echo "Running: $composeCommand"
    $composeCommand || (app_echo "Failed to up the compose stack" && return 1)
    return $?
}

app_compose_down() {
    local composeFile=$1

    if [ -z "$composeFile" ]; then
        app_echo_red "call requires a composeFile"
        return 1
    fi

    composeCommand="docker compose -f $composeFile down --timeout=0 --volumes --remove-orphans"
    app_echo "Running: $composeCommand"
    $composeCommand || ( app_echo "Failed to down the compose stack" && return 1 )
    return $?
}
