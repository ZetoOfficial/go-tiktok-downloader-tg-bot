# 🕺 Telegram Media Downloader Bot

Телеграм-бот для скачивания видео с TikTok и YouTube Shorts.  
Скачивание выполняется напрямую с помощью `yt-dlp`.

---

### 📦 Стек

- Go (Golang)
- Telegram Bot API (`tgbotapi`)
- Docker / Docker Compose
- `yt-dlp`
- FFmpeg

---

## 🚀 Быстрый старт

### 🔧 Настройка `config.yaml`

Пример:

```yaml
bot_token: "TOKEN"
```

---

### 🐳 Запуск через Docker Compose

```bash
docker-compose up --build
```

Это поднимет Телеграм-бот (`tiktok-bot`).

---

### ✅ Примеры команд для пользователя

Отправьте боту ссылку на видео с TikTok — он вернёт видеофайл без водяных знаков.

Пример:
```
https://vt.tiktok.com/ZSMoXTxvS/
```

---

## 🛠 Сборка и запуск вручную (без Docker)

```bash
go build -o bot ./cmd/bot
./bot --config=config/config.yaml
```

---

## 🧼 Полезные команды

```bash
# Перезапустить и пересобрать всё
docker-compose down --volumes --remove-orphans
docker-compose up --build

# Открыть терминал в контейнере бота
docker exec -it tiktok-bot sh
```
