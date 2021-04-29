#!/bin/sh

VSN=1.0.0

curl -LO https://github.com/ipinfo/cli/releases/download/grepip-${VSN}/grepip_${VSN}.deb
sudo dpkg -i grepip_${VSN}.deb
rm grepip_${VSN}.deb

echo
echo 'You can now run `grepip`'.

rm $0
