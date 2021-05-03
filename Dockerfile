FROM golang:1.16.3-alpine
RUN apk update
WORKDIR /goAuth
COPY go.mod go.sum /
RUN go install ./...