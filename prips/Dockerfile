# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /prips
COPY build/prips ./
RUN ls -la
ENTRYPOINT ["/prips/prips"]
