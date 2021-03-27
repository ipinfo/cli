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
$ROOT/scripts/build-all-platforms.sh "$VSN"

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
rm -rf $ROOT/dist/usr
mkdir -p $ROOT/dist/usr/local/bin
cp $ROOT/build/ipinfo_1.0.0b2_linux_amd64 $ROOT/dist/usr/local/bin/ipinfo
dpkg-deb --build ${ROOT}/dist build/ipinfo_${VSN}.deb

# release
gh release create $VSN                                                        \
    -R ipinfo/cli                                                             \
    $ROOT/build/*.tar.gz                                                      \
    $ROOT/build/*.zip                                                         \
    $ROOT/build/*.deb
