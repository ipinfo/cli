#!/bin/sh

VSN=1.0.0

curl -LO https://github.com/ipinfo/cli/releases/download/matchip-${VSN}/matchip${VSN}.deb
sudo dpkg -i matchip_${VSN}.deb
rm matchip_${VSN}.deb

echo
echo 'You can now run `matchip`'.

if [ -f "$0" ]; then
    rm $0
fi
