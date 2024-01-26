#!/bin/bash

# Build binary for all platforms for version $1.
# Optional param LINUX_ONLY can be set to `true`, to build for linux only.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

VSN=$1
LINUX_ONLY=$2

$ROOT/scripts/build-all-platforms.sh "ipinfo" $VSN $LINUX_ONLY
