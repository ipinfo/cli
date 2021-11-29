# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /grepip
COPY build/grepip ./
RUN ls -la
ENTRYPOINT ["/grepip/grepip"]
