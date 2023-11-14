# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /grepdom
COPY build/grepdom ./
ENTRYPOINT ["/grepdom/grepdom"]
