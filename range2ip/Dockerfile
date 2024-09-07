# syntax=docker/dockerfile:1
FROM alpine:3.20.3

WORKDIR /range2ip
COPY build/range2ip ./
ENTRYPOINT ["/range2ip/range2ip"]
