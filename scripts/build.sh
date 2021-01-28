#!/bin/bash

# Build local binary.

DIR=`dirname $0`
ROOT=$DIR/..

go build                                                                      \
    -o $ROOT/build/ipinfo                                                     \
    $ROOT/ipinfo
