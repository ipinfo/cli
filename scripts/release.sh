#!/bin/bash

# Build and upload (to GitHub) for all platforms for version $1.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

VSN=$1

if [ -z "$VSN" ]; then
    echo "require version as first parameter" 2>&1
    exit 1
fi

# build
rm ./build/ipinfo_${VSN}_*
$ROOT/scripts/build-all-platforms.sh "$VSN"

# archive
cd ./build
for t in ./ipinfo_${VSN}_* ; do
    if [[ $t == ./ipinfo_*_windows_* ]]; then
        zip -q ${t/.exe/.zip} $t
    else
        tar -czf ${t}.tar.gz $t
    fi
done
cd ..

# release
gh release create $VSN                                                        \
    ./build/*.tar.gz                                                          \
    ./build/*.zip
