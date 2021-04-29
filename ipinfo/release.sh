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
rm -f $ROOT/build/ipinfo_${VSN}*
$ROOT/ipinfo/build-all-platforms.sh "$VSN"

# archive
cd $ROOT/build
for t in ipinfo_${VSN}_* ; do
    if [[ $t == ipinfo_*_windows_* ]]; then
        zip -q ${t/.exe/.zip} $t
    else
        tar -czf ${t}.tar.gz $t
    fi
done
cd ..

# dist: debian
rm -rf $ROOT/ipinfo/dist/usr
mkdir -p $ROOT/ipinfo/dist/usr/local/bin
cp $ROOT/build/ipinfo_${VSN}_linux_amd64 $ROOT/ipinfo/dist/usr/local/bin/ipinfo
dpkg-deb --build ${ROOT}/ipinfo/dist build/ipinfo_${VSN}.deb

# release
gh release create ipinfo-${VSN}                                               \
    -R ipinfo/cli                                                             \
    -t "ipinfo-${VSN}"                                                        \
    $ROOT/build/ipinfo_*.tar.gz                                               \
    $ROOT/build/ipinfo_*.zip                                                  \
    $ROOT/build/ipinfo_*.deb                                                  \
    $ROOT/ipinfo/macos.sh
