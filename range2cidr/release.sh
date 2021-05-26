#!/bin/bash

# Build and upload (to GitHub) for all platforms for version $1.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

VSN=$1

$ROOT/scripts/release.sh "range2cidr" $VSN
