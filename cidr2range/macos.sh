#!/bin/sh

VSN=1.1.1
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/cidr2range-${VSN}/cidr2range_${VSN}_${PLAT}.tar.gz
tar -xf cidr2range_${VSN}_${PLAT}.tar.gz
rm cidr2range_${VSN}_${PLAT}.tar.gz
mv cidr2range_${VSN}_${PLAT} /usr/local/bin/cidr2range

echo
echo 'You can now run `cidr2range`'.

if [ -f "$0" ]; then
    rm $0
fi
