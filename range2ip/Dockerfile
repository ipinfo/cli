# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /range2ip
COPY build/range2ip ./
RUN ls -la
ENTRYPOINT ["/range2ip/range2ip"]
