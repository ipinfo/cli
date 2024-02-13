#!/bin/sh

if [ "$(id -u)" -ne 0 ]; then
    echo "This script requires root privileges. Please run it as root or with sudo." >&2
    exit 1
fi

VSN=1.2.3
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
curl -LO https://github.com/ipinfo/cli/releases/download/grepip-${VSN}/grepip_${VSN}_linux_${ARCH_NAME}.deb
dpkg -i grepip_${VSN}_linux_${ARCH_NAME}.deb
rm grepip_${VSN}_linux_${ARCH_NAME}.deb

echo
echo 'You can now run `grepip`'.

if [ -f "$0" ]; then
    rm $0
fi
