#!/bin/bash

# Build local binary.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

$ROOT/scripts/build.sh "randip"
