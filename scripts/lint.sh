#!/bin/bash

# Lint all files in the project.

golint                                                                        \
    ./lib/...                                                                 \
    ./ipinfo/...                                                              \
    ./grepip/...                                                              \
    ./cidr2range/...                                                          \
    ./range2cidr/...
    ./range2ip/...
