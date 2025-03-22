# 🕺 Telegram TikTok/Douyin Downloader Bot

Телеграм-бот для скачивания видео с TikTok и Douyin (抖音), без водяных знаков.  
Использует внешний API: [`evil0ctal/douyin_tiktok_download_api`](https://github.com/Evil0ctal/Douyin_TikTok_Download_API)

---

### 📦 Стек

- Go (Golang)
- Telegram Bot API (`tgbotapi`)
- Docker / Docker Compose
- Douyin API (через внешний Docker-образ)

---

## 🚀 Быстрый старт

### 🔧 Настройка `config.yaml`

Пример:

```yaml
bot_token: "TOKEN"
douyin_api_url: "http://douyin-api:80"
```

---

### 🐳 Запуск через Docker Compose

```bash
docker-compose up --build
```

Это поднимет:

- Телеграм-бот (`tiktok-bot`)
- Внешний Douyin API (`douyin-api`)

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
