#!/bin/bash

# Builds cli $1 version $2 for all platforms and packages them for release.
# Optional param LINUX_ONLY can be set to `true`, to build for linux only.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

CLI=$1
VSN=$2
LINUX_ONLY=$3

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
rm -f $ROOT/build/${CLI}_${VSN}*
$ROOT/${CLI}/build-all-platforms.sh "$VSN" "$LINUX_ONLY"

# archive
cd $ROOT/build
for t in ${CLI}_${VSN}_* ; do
    if [[ $t == ${CLI}_*_windows_* ]]; then
        zip -q ${t/.exe/.zip} $t
    else
        tar -czf ${t}.tar.gz $t
    fi
done
cd ..

# dist: debian
declare -A debian_archs
debian_archs[linux_386]="i386"
debian_archs[linux_amd64]="amd64"
debian_archs[linux_arm]="armhf"
debian_archs[linux_arm64]="arm64"
for t in                                                                      \
    linux_386                                                                 \
    linux_amd64                                                               \
    linux_arm                                                                 \
    linux_arm64;
do
    os="${t%_*}"
    arch="${t#*_}"
    debian_arch="${debian_archs[$t]}"
    output="${CLI}_${VSN}_${os}_${arch}"

    # Update Version and Architecture in the control file
    sed -i "s/Version: .*/Version: $VSN/" "${ROOT}/${CLI}/dist/DEBIAN/control"
    sed -i "s/Architecture: .*/Architecture: $debian_arch/" "${ROOT}/${CLI}/dist/DEBIAN/control"

    rm -rf "$ROOT/${CLI}/dist/usr"
    mkdir -p "$ROOT/${CLI}/dist/usr/local/bin"
    cp "$ROOT/build/${CLI}_${VSN}_${os}_${arch}" "$ROOT/${CLI}/dist/usr/local/bin/${CLI}"
    dpkg-deb -Zgzip --build "${ROOT}/${CLI}/dist" "build/${output}.deb"
done

wait
