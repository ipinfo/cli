#!/bin/bash

# Build and upload (to Dockerhub) for all platforms for version $1.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

VSN=$1
RELEASE=$2

$ROOT/scripts/docker.sh "cidr2range" $VSN $RELEASE
