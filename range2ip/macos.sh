#!/bin/sh

VSN=1.2.0
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/range2ip-${VSN}/range2ip_${VSN}_${PLAT}.tar.gz
tar -xf range2ip_${VSN}_${PLAT}.tar.gz
rm range2ip_${VSN}_${PLAT}.tar.gz
mv range2ip_${VSN}_${PLAT} /usr/local/bin/range2ip

echo
echo 'You can now run `range2ip`'.

if [ -f "$0" ]; then
    rm $0
fi
