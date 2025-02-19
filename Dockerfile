FROM golang:1.23.0-alpine AS builder
LABEL authors="guilh"

RUN apk add --no-cache make gcc libc-dev

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY cert.pem ./
COPY key.pem ./
RUN go mod tidy

COPY . .

RUN go build -o main main.go

RUN chmod +x main

RUN whoami
RUN id

EXPOSE 7890

CMD ["./main"]