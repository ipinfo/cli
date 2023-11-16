#!/bin/sh

VSN=1.0.0

curl -LO https://github.com/ipinfo/cli/releases/download/grepdom-${VSN}/grepdom_${VSN}.deb
sudo dpkg -i grepdom_${VSN}.deb
rm grepdom_${VSN}.deb

echo
echo 'You can now run `grepdom`'.

if [ -f "$0" ]; then
    rm $0
fi
