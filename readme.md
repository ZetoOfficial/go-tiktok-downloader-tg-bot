# 🕺 Telegram Media Downloader Bot

Телеграм-бот для скачивания видео с TikTok, YouTube Shorts и публичных Instagram Reels.
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

Отправьте боту ссылку на TikTok, YouTube Shorts или публичный Instagram Reel — он вернёт видеофайл.

Примеры:
```
https://vt.tiktok.com/ZSMoXTxvS/
https://www.youtube.com/shorts/VIDEO_ID
https://www.instagram.com/reel/REEL_ID/
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

---

## 🚢 Деплой (CI/CD)

Пайплайн — `.github/workflows/ci.yml`:

1. **checks** — `golangci-lint` + `go test ./...` на любой ветке.
2. **docker-build** — проверка, что образ собирается (только не на `master`, никуда не пушит).
3. **build-push** — только на `master`: собирает образ и пушит в GHCR с тегами `sha-<short>` и `latest`.
4. **deploy** — копирует `docker-compose.prod.yml` и `deploy.sh` на сервер по `scp`, затем
   `docker compose pull && up -d` в `/opt/tiktok-bot`.

На сервере в `/opt/tiktok-bot/` лежат: `docker-compose.prod.yml`, `deploy.sh` (оба
перезаписываются CI на каждый деплой) и `config.yaml` с токеном бота — **только вручную**,
CI его не трогает.

**Откат:** GitHub → Actions → CI → Run workflow → ветка `master` → в `image_tag` указать
`sha-<short>` нужного коммита. Текущий выкаченный тег лежит в `/opt/tiktok-bot/.deployed-tag`.

**Секреты GitHub:** `DEPLOY_HOST`, `DEPLOY_USER`, `DEPLOY_SSH_KEY` (приватный ключ),
`SSH_KNOWN_HOSTS` (вывод `ssh-keyscan`).
