FROM golang:1.23.0-alpine AS builder
LABEL authors="guilh"

# Instalar dependências para compilar o Go
RUN apk add --no-cache make gcc libc-dev

WORKDIR /app

# Copiar os arquivos go.mod, go.sum e código da aplicação para o contêiner
COPY go.mod ./
COPY go.sum ./
COPY . ./

# Instalar dependências Go
RUN go mod tidy

# Construir o binário
RUN go build -o main main.go

# Tornar o binário executável
RUN chmod +x main

# Copiar os certificados SSL para o contêiner (ajuste os caminhos conforme necessário)
COPY /etc/letsencrypt/live/seu-dominio.com.br/cert.pem /app/cert.pem
COPY /etc/letsencrypt/live/seu-dominio.com.br/privkey.pem /app/key.pem

# Expor a porta que o servidor vai escutar
EXPOSE 7890

# Definir o comando de entrada para rodar o binário
CMD ["./main"]
