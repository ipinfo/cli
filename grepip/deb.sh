#!/bin/sh

VSN=1.0.3

curl -LO https://github.com/ipinfo/cli/releases/download/grepip-${VSN}/grepip_${VSN}.deb
sudo dpkg -i grepip_${VSN}.deb
rm grepip_${VSN}.deb

echo
echo 'You can now run `grepip`'.

if [ -f "$0" ]; then
    rm $0
fi
