#!/bin/bash

# Regenerate ctags.

ctags \
    --recurse=yes \
    --exclude=node_modules \
    --exclude=dist \
    --exclude=build \
    --exclude=target \
    -f .vim/tags \
    --tag-relative=never \
    --totals=yes \
        ./lib \
        ./ipinfo \
        ./grepdomain \
        ./grepip \
        ./matchip
