#!/bin/sh

VSN=1.0.0

curl -LO https://github.com/ipinfo/cli/releases/download/prips-${VSN}/prips_${VSN}.deb
sudo dpkg -i prips_${VSN}.deb
rm prips_${VSN}.deb

echo
echo 'You can now run `prips`'.

if [ -f "$0" ]; then
    rm $0
fi
