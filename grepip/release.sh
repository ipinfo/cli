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
rm -f $ROOT/build/grepip_${VSN}*
$ROOT/grepip/build-all-platforms.sh "$VSN"

# archive
cd $ROOT/build
for t in grepip_${VSN}_* ; do
    if [[ $t == grepip_*_windows_* ]]; then
        zip -q ${t/.exe/.zip} $t
    else
        tar -czf ${t}.tar.gz $t
    fi
done
cd ..

# dist: debian
rm -rf $ROOT/grepip/dist/usr
mkdir -p $ROOT/grepip/dist/usr/local/bin
cp $ROOT/build/grepip_${VSN}_linux_amd64 $ROOT/grepip/dist/usr/local/bin/grepip
dpkg-deb --build ${ROOT}/grepip/dist build/grepip_${VSN}.deb

# release
gh release create grepip-${VSN}                                               \
    -R ipinfo/cli                                                             \
    -t "grepip-${VSN}"                                                        \
    $ROOT/build/grepip_*.tar.gz                                               \
    $ROOT/build/grepip_*.zip                                                  \
    $ROOT/build/grepip_*.deb
