#!/bin/bash

# Build binary for cli $1.

DIR=`dirname $0`
ROOT=$DIR/..

CLI=$1

go build                                                                      \
    -o $ROOT/build/${CLI}                                                     \
    $ROOT/${CLI}
