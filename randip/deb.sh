#!/bin/sh

VSN=1.0.0

curl -LO https://github.com/ipinfo/cli/releases/download/randip-${VSN}/randip_${VSN}.deb
sudo dpkg -i randip_${VSN}.deb
rm randip_${VSN}.deb

echo
echo 'You can now run `randip`'.

if [ -f "$0" ]; then
    rm $0
fi
