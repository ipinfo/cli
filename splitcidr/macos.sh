#!/bin/sh

VSN=1.0.0
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/splitcidr-${VSN}/splitcidr_${VSN}_${PLAT}.tar.gz && \
tar -xf splitcidr_${VSN}_${PLAT}.tar.gz && \
rm splitcidr_${VSN}_${PLAT}.tar.gz && \
sudo mv splitcidr_${VSN}_${PLAT} /usr/local/bin/splitcidr && \
echo 'You can now run `splitcidr`.'

if [ -f "$0" ]; then
    rm "$0"
fi
