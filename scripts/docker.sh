#!/bin/bash

# Build (and optionally upload to Dockerhub) for CLI $1 & version $2.
# version ($2) will be set to `latest` in case not provided.
# additionally --release, -r to release container to Dockerhub.
# note: if you want to release as well, you should be logged in with `docker`.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

CLI=$1
VSN=$2
RELEASE=$3

if [ -z "$CLI" ]; then
    echo "require cli as first parameter" 2>&1
    exit 1
fi

if [ -z "$VSN" ]; then
    VSN="latest"
fi


# build
# disable CGO to make static binary
CGO_ENABLED=0 go build                                                        \
    -o $ROOT/${CLI}/build/$CLI                                                \
    $ROOT/${CLI}

# docker container
docker build --tag ipinfo/$CLI:$VSN $ROOT/$CLI/

# cleanup 
rm -r $ROOT/$CLI/build

if [ "$RELEASE" = "-r" ] || [ "$RELEASE" = "--release" ]; then
    # push on docker hub
    docker push ipinfo/$CLI:$VSN
fi
