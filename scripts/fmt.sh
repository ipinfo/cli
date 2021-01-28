#!/bin/bash

DIR=`dirname $0`
ROOT=$DIR/..

# Format code in project.

gofmt -w \
    $ROOT/ipinfo
