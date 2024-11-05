FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/app ./cmd/app/main.go

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/app .

RUN chmod +x app

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

EXPOSE 9999
EXPOSE 5173

CMD ["./app"]
