FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bot ./cmd/bot

FROM denoland/deno:alpine

WORKDIR /app

# yt-dlp с EJS-компонентами и Deno для обработки JavaScript-проверок YouTube.
# ffmpeg склеивает раздельные видео/аудио дорожки (например, YouTube Shorts).
RUN apk add --no-cache python3 py3-pip ffmpeg ca-certificates \
    && pip install --no-cache-dir --break-system-packages "yt-dlp[default]"

COPY --from=builder /app/bot .

COPY config/config.yaml ./config.yaml

CMD ["./bot", "--config=./config.yaml"]
