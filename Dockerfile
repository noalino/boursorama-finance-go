FROM golang:1.15-alpine

WORKDIR /go/src/boursorama-finance-go
COPY . .

RUN go install -v ./...
