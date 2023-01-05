# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /ipinfo
COPY build/ipinfo ./
ENTRYPOINT ["/ipinfo/ipinfo"]
