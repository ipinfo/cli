#!/bin/bash

DIR=`dirname $0`
ROOT=$DIR/..

# Format code in project.

gofmt -w                                                                      \
    $ROOT/lib                                                                 \
    $ROOT/ipinfo                                                              \
    $ROOT/grepip                                                              \
    $ROOT/cidr2range                                                          \
    $ROOT/range2cidr                                                          \
    $ROOT/range2ip                                                            
