# syntax=docker/dockerfile:1
FROM alpine:latest

WORKDIR /grepdomain
COPY build/grepdomain ./
ENTRYPOINT ["/grepdomain/grepdomain"]
