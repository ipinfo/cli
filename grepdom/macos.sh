#!/bin/sh

VSN=1.0.0
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/grepdom-${VSN}/grepdom_${VSN}_${PLAT}.tar.gz
tar -xf grepdom_${VSN}_${PLAT}.tar.gz
rm grepdom_${VSN}_${PLAT}.tar.gz
mv grepdom_${VSN}_${PLAT} /usr/local/bin/grepdom

echo
echo 'You can now run `grepdom`'.

if [ -f "$0" ]; then
    rm $0
fi
