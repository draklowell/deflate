FROM golang:1.24.0 AS builder-image

WORKDIR /builder
COPY . .

RUN mkdir /build
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/deflate .
