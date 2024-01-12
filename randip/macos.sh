#!/bin/sh

VSN=1.1.0
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/randip-${VSN}/randip_${VSN}_${PLAT}.tar.gz && \
tar -xf randip_${VSN}_${PLAT}.tar.gz && \
rm randip_${VSN}_${PLAT}.tar.gz && \
sudo mv randip_${VSN}_${PLAT} /usr/local/bin/randip && \
echo 'You can now run `randip`.'

if [ -f "$0" ]; then
    rm "$0"
fi
