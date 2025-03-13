package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/config"
	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/api"
	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/bot"
	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/clients"
	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "Путь к файлу конфигурации")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Ошибка инициализации Telegram Bot: %v", err)
	}
	botAPI.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := botAPI.GetUpdatesChan(u)

	douyinClient, err := clients.NewDouyinClient(cfg)
	if err != nil {
		log.Fatalf("Ошибка инициализации DouyinClient: %v", err)
	}

	sender := api.NewTelegramAdapter(botAPI)

	downloaderService := service.NewDownloadService(douyinClient)
	messageService := service.NewMessageService(sender)

	handler := bot.NewHandler(downloaderService, messageService)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{})

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case update, ok := <-updates:
				if !ok {
					return
				}
				go handler.HandleUpdate(botAPI, update)
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		<-sigs
		log.Println("Выключение бота...")
		cancel()
		done <- struct{}{}
	}()

	<-done
	log.Println("Бот успешно завершил работу")
}
