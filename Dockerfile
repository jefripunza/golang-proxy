FROM oven/bun:latest AS fe-builder
WORKDIR /app
COPY package.json bun.lock ./
RUN bun install
COPY . .
RUN bun run build

FROM golang:1.24-bookworm AS be-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=fe-builder /app/dist ./dist
RUN go build -o go-proxy .

FROM debian:bookworm-slim
LABEL org.opencontainers.image.authors="Jefri Herdi Triyanto <jefriherditriyanto@gmail.com>"
LABEL description="Golang Proxy — dynamic reverse proxy with web dashboard, SSL, and IP blocklist."

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /app/database

WORKDIR /app
COPY --from=be-builder /app/go-proxy .

EXPOSE 80 443 8080 8000
CMD ["./go-proxy"]
