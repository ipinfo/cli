#!/bin/sh

VSN=1.0.0

curl -LO https://github.com/ipinfo/cli/releases/download/cidr2range-${VSN}/cidr2range_${VSN}.deb
sudo dpkg -i cidr2range_${VSN}.deb
rm cidr2range_${VSN}.deb

echo
echo 'You can now run `cidr2range`'.

if [ -f "$0" ]; then
    rm $0
fi
