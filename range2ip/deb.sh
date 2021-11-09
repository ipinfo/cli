#!/bin/sh

VSN=1.0.0

curl -LO https://github.com/ipinfo/cli/releases/download/range2ip-${VSN}/range2ip_${VSN}.deb
sudo dpkg -i range2ip_${VSN}.deb
rm range2ip_${VSN}.deb

echo
echo 'You can now run `range2ip`'.

if [ -f "$0" ]; then
    rm $0
fi
