FROM golang:1.15-alpine

WORKDIR /go/src/go-fetch-quotes
COPY . .

RUN go install -v ./...
