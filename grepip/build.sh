#!/bin/bash

# Build local binary.

DIR=`dirname $0`
ROOT=$DIR/..

go build                                                                      \
    -o $ROOT/build/grepip                                                     \
    $ROOT/grepip
