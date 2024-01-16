#!/bin/sh

set -e

VSN=1.2.3
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/grepip-${VSN}/grepip_${VSN}_${PLAT}.tar.gz
tar -xf grepip_${VSN}_${PLAT}.tar.gz
rm grepip_${VSN}_${PLAT}.tar.gz
sudo mv grepip_${VSN}_${PLAT} /usr/local/bin/grepip

echo
echo 'You can now run `grepip`.'

if [ -f "$0" ]; then
    rm "$0"
fi
