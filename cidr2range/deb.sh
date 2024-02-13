#!/bin/sh

if [ "$(id -u)" -ne 0 ]; then
    echo "This script requires root privileges. Please run it as root or with sudo." >&2
    exit 1
fi

VSN=1.2.0
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
curl -LO https://github.com/ipinfo/cli/releases/download/cidr2range-${VSN}/cidr2range_${VSN}_linux_${ARCH_NAME}.deb
dpkg -i cidr2range_${VSN}_linux_${ARCH_NAME}.deb
rm cidr2range_${VSN}_linux_${ARCH_NAME}.deb

echo
echo 'You can now run `cidr2range`'.

if [ -f "$0" ]; then
    rm $0
fi
