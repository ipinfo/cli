#!/bin/sh

VSN=2.3.1
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/ipinfo-${VSN}/ipinfo_${VSN}_${PLAT}.tar.gz
tar -xf ipinfo_${VSN}_${PLAT}.tar.gz
rm ipinfo_${VSN}_${PLAT}.tar.gz
mv ipinfo_${VSN}_${PLAT} /usr/local/bin/ipinfo

echo
echo 'You can now run `ipinfo`'.

if [ -f "$0" ]; then
    rm $0
fi
