FROM golang:1.16.3-alpine
RUN apk update && apk add git 
WORKDIR /goAuth