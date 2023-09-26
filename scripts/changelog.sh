#!/bin/sh

# NOTE: Github actions uses this script.
# Get changelog of cli $1 & version $2 from CHANGELOG.md.

set -e

DIR=`dirname $0`
ROOT=$DIR/..

CLI=$1
VERSION=$2

MARKER_PREFIX="##"
found=0

cat $ROOT/${CLI}/CHANGELOG.md | while read "line"; do

    # Find the version heading
    if [ $found -eq 0 ] && echo "$line" | grep -q "^$MARKER_PREFIX $VERSION$"; then
        found=1
        continue
    fi

    # Reaching next delimter - stop
    if [ $found -eq 1 ] && echo "$line" | grep -q -E "^$MARKER_PREFIX [[:digit:]]+\.[[:digit:]]+\.[[:digit:]]+"; then
        found=0
        break
    fi

    # Keep printing out lines as no other version delimiter found
    if [ $found -eq 1 ]; then
        echo "$line"
    fi
done
