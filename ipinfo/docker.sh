#!/bin/bash

# Build and upload (to Dockerhub) for all platforms for version $1.
# Use `-r` or `--release` after the version to also push to Dockerhub.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

VSN=$1
RELEASE=$2

$ROOT/scripts/docker.sh "ipinfo" $VSN $RELEASE
