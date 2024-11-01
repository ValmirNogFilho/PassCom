# Usar uma imagem base do Go
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app cmd/app/main.go
RUN go build -o api cmd/api/api.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/api .

RUN chmod +x app api
RUN apk add --no-cache libc6-compat

EXPOSE 9999
EXPOSE 8888
