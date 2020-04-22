FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Ruben Eguiluz <regiluze@gmail.com>"

RUN mkdir -p /go/src/github.com/regiluze/form3-account-api-client

COPY . /go/src/github.com/regiluze/form3-account-api-client

WORKDIR /go/src/github.com/regiluze/form3-account-api-client
RUN make deps