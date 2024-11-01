# Usar uma imagem base do Go
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilar o binário com caminho absoluto para evitar problemas de cópia
RUN go build -o /app/app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app

# Copiar o binário do estágio de construção
COPY --from=builder /app/app .

# Tornar o binário executável e instalar compatibilidade com libc
RUN chmod +x app
RUN apk add --no-cache libc6-compat

EXPOSE 8080

CMD ["./app"]
