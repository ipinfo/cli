#!/bin/bash

# Build and upload (to GitHub) for all platforms for cli $1 & version $2.

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
    echo "\"$CLI\" is not a cli" 2>&1
    exit 1
fi

# build
$ROOT/scripts/build-archive-all.sh "$CLI" "$VSN"

# release
gh release create ${CLI}-${VSN}                                               \
    -R ipinfo/cli                                                             \
    -t "${CLI}-${VSN}"                                                        \
    $ROOT/build/${CLI}_${VSN}*.tar.gz                                         \
    $ROOT/build/${CLI}_${VSN}*.zip                                            \
    $ROOT/build/${CLI}_${VSN}*.deb                                            \
    $ROOT/${CLI}/macos.sh                                                     \
    $ROOT/${CLI}/windows.ps1                                                  \
    $ROOT/${CLI}/deb.sh
