FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bot ./cmd/bot

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bot .

COPY config/config.yaml ./config.yaml

CMD ["./bot", "--config=./config.yaml"]
