# Usar uma imagem base do Go para a etapa de build
FROM golang:1.22 AS builder

WORKDIR /app

# Copiar os arquivos go.mod e go.sum e baixar dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código-fonte
COPY . .

# Compilar o binário
RUN go build -o /app/app ./cmd/app/main.go

# Imagem final usando Alpine para manter a imagem leve
FROM alpine:latest

WORKDIR /app

# Copiar o binário compilado da etapa de build
COPY --from=builder /app/app .

# Tornar o binário executável e instalar compatibilidade com libc para bibliotecas Go
RUN chmod +x app
RUN apk add --no-cache libc6-compat

# Expor a porta que a aplicação vai usar
EXPOSE 8080

# Comando para iniciar o aplicativo
CMD ["./app"]
