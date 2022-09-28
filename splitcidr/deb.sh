#!/bin/sh

VSN=1.0.0

curl -LO https://github.com/ipinfo/cli/releases/download/splitcidr-${VSN}/splitcidr_${VSN}.deb
sudo dpkg -i splitcidr_${VSN}.deb
rm splitcidr_${VSN}.deb

echo
echo 'You can now run `splitcidr`'.

if [ -f "$0" ]; then
    rm $0
fi
