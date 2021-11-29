# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /cidr2ip
COPY build/cidr2ip ./
ENTRYPOINT ["/cidr2ip/cidr2ip"]
