FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bot ./cmd/bot

FROM alpine:latest

WORKDIR /app

# yt-dlp (через pip — свежая версия, важно чтобы не отставать от TikTok/YouTube)
# + ffmpeg для склейки видео/аудио дорожек (например, YouTube Shorts).
RUN apk add --no-cache python3 py3-pip ffmpeg ca-certificates \
    && pip install --no-cache-dir --break-system-packages yt-dlp

COPY --from=builder /app/bot .

COPY config/config.yaml ./config.yaml

CMD ["./bot", "--config=./config.yaml"]
