#!/bin/sh

set -e

VSN=1.0.0
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/matchip-${VSN}/matchip_${VSN}_${PLAT}.tar.gz
tar -xf matchip_${VSN}_${PLAT}.tar.gz
rm matchip_${VSN}_${PLAT}.tar.gz
sudo mv matchip_${VSN}_${PLAT} /usr/local/bin/matchip

echo
echo 'You can now run `matchip`.'

if [ -f "$0" ]; then
    rm "$0"
fi
