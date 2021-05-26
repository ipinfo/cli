#!/bin/sh

VSN=1.0.0
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/range2cidr-${VSN}/range2cidr_${VSN}_${PLAT}.tar.gz
tar -xf range2cidr_${VSN}_${PLAT}.tar.gz
rm range2cidr_${VSN}_${PLAT}.tar.gz
mv range2cidr_${VSN}_${PLAT} /usr/local/bin/range2cidr

echo
echo 'You can now run `range2cidr`'.

if [ -f "$0" ]; then
    rm $0
fi
