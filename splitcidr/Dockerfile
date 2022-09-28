# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /splitcidr
COPY build/splitcidr ./
ENTRYPOINT ["/splitcidr/splitcidr"]
