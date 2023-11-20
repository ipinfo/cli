#!/bin/sh

VSN=1.0.0
PLAT=darwin_amd64

curl -LO https://github.com/ipinfo/cli/releases/download/grepdomain-${VSN}/grepdomain_${VSN}_${PLAT}.tar.gz
tar -xf grepdomain_${VSN}_${PLAT}.tar.gz
rm grepdomain_${VSN}_${PLAT}.tar.gz
mv grepdomain_${VSN}_${PLAT} /usr/local/bin/grepdomain

echo
echo 'You can now run `grepdomain`'.

if [ -f "$0" ]; then
    rm $0
fi
