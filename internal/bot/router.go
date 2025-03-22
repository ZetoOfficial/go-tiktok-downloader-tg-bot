package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleCommand(botAPI *tgbotapi.BotAPI, message *tgbotapi.Message) {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		{Command: "help", Description: "Справка по командам"},
	}

	_, err := botAPI.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		log.Fatalf("Ошибка при установке команд: %v", err)
	}

	chatID := message.Chat.ID
	switch message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(chatID, "Добро пожаловать! Отправьте ссылку на видео или фото из TikTok.")
		_, _ = botAPI.Send(msg)
	case "help":
		helpText := "Функционал бота:\n" +
			"- Отправьте ссылку на TikTok для загрузки видео/изображений.\n" +
			"- Если видео больше 50MB, оно не будет отправлено."
		msg := tgbotapi.NewMessage(chatID, helpText)
		_, _ = botAPI.Send(msg)
	default:
		msg := tgbotapi.NewMessage(chatID, "Неизвестная команда.")
		_, _ = botAPI.Send(msg)
	}
}
