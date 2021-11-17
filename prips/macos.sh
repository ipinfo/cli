#!/bin/sh

VSN=1.0.0
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/prips-${VSN}/prips_${VSN}_${PLAT}.tar.gz
tar -xf prips_${VSN}_${PLAT}.tar.gz
rm prips_${VSN}_${PLAT}.tar.gz
mv prips_${VSN}_${PLAT} /usr/local/bin/prips

echo
echo 'You can now run `prips`'.

if [ -f "$0" ]; then
    rm $0
fi
