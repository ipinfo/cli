# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /cidr2range
COPY build/cidr2range ./
ENTRYPOINT ["/cidr2range/cidr2range"]
