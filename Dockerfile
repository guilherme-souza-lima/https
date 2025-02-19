# Usando uma imagem base de Go
FROM golang:1.23-alpine

# Definindo o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiando os arquivos do projeto para dentro do contêiner
COPY . .

# Expondo a porta 7890
EXPOSE 7890

# Comando para rodar a aplicação
CMD ["go", "run", "main.go"]
