#!/bin/bash

# Build binary for all platforms for version $1.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

VSN=$1

if [ -z "$VSN" ]; then
    echo "require version as first parameter" 2>&1
    exit 1
fi

for t in                                                                      \
    darwin_amd64                                                              \
    dragonfly_amd64                                                           \
    freebsd_amd64                                                             \
    linux_amd64                                                               \
    netbsd_amd64                                                              \
    openbsd_amd64                                                             \
    plan9_amd64                                                               \
    solaris_amd64                                                             \
    windows_amd64 ;
do
    os="${t%_*}"
    arch="${t#*_}"
    output="ipinfo_${VSN}_${os}_${arch}"

    if [ "$os" == "windows" ] ; then
        output+=".exe"
    fi

    GOOS=$os GOARCH=$arch go build                                            \
        -o $ROOT/build/${output}                                              \
        $ROOT/ipinfo &
done

wait
