#!/bin/sh

VSN=1.0.0

curl -LO https://github.com/ipinfo/cli/releases/download/cidr2ip-${VSN}/cidr2ip_${VSN}.deb
sudo dpkg -i cidr2ip_${VSN}.deb
rm cidr2ip_${VSN}.deb

echo
echo 'You can now run `cidr2ip`'.

if [ -f "$0" ]; then
    rm $0
fi
