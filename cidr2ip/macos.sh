#!/bin/sh

set -e

VSN=1.0.0
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/cidr2ip-${VSN}/cidr2ip_${VSN}_${PLAT}.tar.gz
tar -xf cidr2ip_${VSN}_${PLAT}.tar.gz
rm cidr2ip_${VSN}_${PLAT}.tar.gz
sudo mv cidr2ip_${VSN}_${PLAT} /usr/local/bin/cidr2ip

echo
echo 'You can now run `cidr2ip`.'

if [ -f "$0" ]; then
    rm "$0"
fi
