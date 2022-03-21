#!/bin/sh

VSN=1.3.0

curl -LO https://github.com/ipinfo/cli/releases/download/range2cidr-${VSN}/range2cidr_${VSN}.deb
sudo dpkg -i range2cidr_${VSN}.deb
rm range2cidr_${VSN}.deb

echo
echo 'You can now run `range2cidr`'.

if [ -f "$0" ]; then
    rm $0
fi
