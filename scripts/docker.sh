#!/bin/bash

# Build and upload (to Dockerhub) for cli $1 & version $2.
# run this script with sudo 
# note: docker should be loggedin in your terminal before running this script

set -e

DIR=`dirname $0`
ROOT=$DIR/..

CLI=$1
VSN=$2

if [ -z "$CLI" ]; then
    echo "require cli as first parameter" 2>&1
    exit 1
fi

if [ -z "$VSN" ]; then
    echo "require version as second parameter" 2>&1
    exit 1
fi


# build
CGO_ENABLED=0 go build                                                        \
    -o $ROOT/${CLI}/build/$CLI                                                \
    $ROOT/${CLI}

# docker container
sudo docker build --tag ipinfo/$CLI:$VSN $ROOT/$CLI/

# push on docker hub
docker push ipinfo/$CLI:$VSN

# cleanup 
sudo rm -rf $ROOT/$CLI/build