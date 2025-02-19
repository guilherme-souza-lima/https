# Etapa de build
FROM golang:1.23 AS builder

WORKDIR /app

# Copia os arquivos do projeto
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila a aplicação
RUN go build -o server .

# Etapa final: imagem mais leve
FROM alpine:latest

WORKDIR /root/

# Instala dependências necessárias
RUN apk --no-cache add ca-certificates

# Copia o binário compilado da etapa anterior
COPY --from=builder /app/server .

# Copia os certificados SSL
COPY cert.pem key.pem ./

# Expõe a porta HTTPS
EXPOSE 7890

# Executa o servidor
CMD ["./server"]
