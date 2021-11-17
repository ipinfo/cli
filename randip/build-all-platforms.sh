#!/bin/bash

# Build binary for all platforms for version $1.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

VSN=$1

$ROOT/scripts/build-all-platforms.sh "randip" $VSN
