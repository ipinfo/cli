# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /randip
COPY build/randip ./
RUN ls -la
ENTRYPOINT ["/randip/randip"]
