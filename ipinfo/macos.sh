#!/bin/sh

VSN=1.1.0
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/ipinfo-${VSN}/ipinfo_${VSN}_${PLAT}.tar.gz
tar -xvf ipinfo_${VSN}_${PLAT}.tar.gz
rm ipinfo_${VSN}_${PLAT}.tar.gz
mv ipinfo_${VSN}_${PLAT} /usr/local/bin/ipinfo
rm $0
