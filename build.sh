#!/bin/bash

for dir in $(ls); do
    [[ -d $dir ]] && [[ -f $dir/dir.go ]] && cd $(pwd)/$dir && go build -o $dir && mv $dir /usr/bin && cd ..
done