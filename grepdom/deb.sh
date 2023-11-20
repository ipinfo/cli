#!/bin/sh

VSN=1.0.0

curl -LO https://github.com/ipinfo/cli/releases/download/grepdomain-${VSN}/grepdomain_${VSN}.deb
sudo dpkg -i grepdomain_${VSN}.deb
rm grepdomain_${VSN}.deb

echo
echo 'You can now run `grepdomain`'.

if [ -f "$0" ]; then
    rm $0
fi
