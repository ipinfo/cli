#!/bin/sh

VSN=3.3.0
DEFAULT_ARCH=amd64

ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH_NAME="amd64"
        ;;
    i386|i686)
        ARCH_NAME="386"
        ;;
    aarch64)
        ARCH_NAME="arm64"
        ;;
    armv7l)
        ARCH_NAME="arm"
        ;;
    *)
        ARCH_NAME=$DEFAULT_ARCH
        ;;
esac

curl -LO https://github.com/ipinfo/cli/releases/download/ipinfo-${VSN}/ipinfo_${VSN}_linux_${ARCH_NAME}.deb
if command -v sudo >/dev/null 2>&1; then
    sudo dpkg -i ipinfo_${VSN}_linux_${ARCH_NAME}.deb
else
    dpkg -i ipinfo_${VSN}_linux_${ARCH_NAME}.deb
fi
rm ipinfo_${VSN}_linux_${ARCH_NAME}.deb

echo
echo 'You can now run `ipinfo`'.

if [ -f "$0" ]; then
    rm $0
fi
