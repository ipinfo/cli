#!/bin/bash

# Build binary for all platforms for cli $1 & version $2.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

CLI=$1
VSN=$2

if [ -z "$CLI" ]; then
    echo "require cli as first parameter" 2>&1
    exit 1
fi

if [ -z "$VSN" ]; then
    echo "require version as second parameter" 2>&1
    exit 1
fi

if [ ! -d "$CLI" ]; then
    echo "\"$cli\" is not a cli" 2>&1
    exit 1
fi

for t in                                                                      \
    darwin_amd64                                                              \
    darwin_arm64                                                              \
    dragonfly_amd64                                                           \
    freebsd_386                                                               \
    freebsd_amd64                                                             \
    freebsd_arm                                                               \
    freebsd_arm64                                                             \
    linux_386                                                                 \
    linux_amd64                                                               \
    linux_arm                                                                 \
    linux_arm64                                                               \
    netbsd_386                                                                \
    netbsd_amd64                                                              \
    netbsd_arm                                                                \
    netbsd_arm64                                                              \
    openbsd_386                                                               \
    openbsd_amd64                                                             \
    openbsd_arm                                                               \
    openbsd_arm64                                                             \
    plan9_386                                                                 \
    plan9_amd64                                                               \
    plan9_arm                                                                 \
    solaris_amd64                                                             \
    windows_386                                                               \
    windows_amd64                                                             \
    windows_arm ;
do
    os="${t%_*}"
    arch="${t#*_}"
    output="${CLI}_${VSN}_${os}_${arch}"

    if [ "$os" == "windows" ] ; then
        output+=".exe"
    fi

    GOOS=$os GOARCH=$arch go build                                            \
        -o $ROOT/build/${output}                                              \
        $ROOT/${CLI} &
done

wait
