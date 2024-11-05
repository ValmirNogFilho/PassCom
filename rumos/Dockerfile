# Usar uma imagem base do Go
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api cmd/api/api.go

FROM debian:bookworm-slim

WORKDIR /app

# Copiar o binário compilado da etapa de build
COPY --from=builder /app/api .

# Tornar o binário executável
RUN chmod +x api

# Instalar dependências mínimas, caso seu app precise
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 7777
