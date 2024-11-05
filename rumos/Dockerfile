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

# Imagem final usando debian-slim para compatibilidade com glibc
FROM debian:bookworm-slim

WORKDIR /app

# Copiar o binário compilado da etapa de build
COPY --from=builder /app/app .

# Tornar o binário executável
RUN chmod +x app

# Instalar dependências mínimas, caso seu app precise
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Expor a porta que a aplicação vai usar
EXPOSE 7777
EXPOSE 5173

# Comando para iniciar o aplicativo
CMD ["./app"]
