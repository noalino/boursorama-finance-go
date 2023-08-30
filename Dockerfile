FROM golang:1.21-alpine

WORKDIR /go/src/boursorama-finance-go
COPY . .

RUN go install -v ./...
