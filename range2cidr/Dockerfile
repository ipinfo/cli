# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /range2cidr
COPY build/range2cidr ./
ENTRYPOINT ["/range2cidr/range2cidr"]
