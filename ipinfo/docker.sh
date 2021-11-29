#!/bin/bash

# Build and upload (to Dockerhub) for all platforms for version $1.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

VSN=$1

$ROOT/scripts/docker.sh "ipinfo" $VSN
