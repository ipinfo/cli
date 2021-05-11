FROM golang:1.16.4-alpine3.13

# install ipinfo-cli
RUN go get github.com/ipinfo/cli/ipinfo