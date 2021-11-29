# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /cidr2range
COPY build/cidr2range ./
RUN ls -la
ENTRYPOINT ["/cidr2range/cidr2range"]
