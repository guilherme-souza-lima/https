# Dockerfile simplificado, sem os comandos COPY para os certificados
FROM golang:1.23.0-alpine AS builder
LABEL authors="guilh"

RUN apk add --no-cache make gcc libc-dev

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main main.go

RUN chmod +x main

EXPOSE 7890

CMD ["./main"]
