package bot

import (
	"context"
	"log"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/sanitizer"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DownloaderService interface {
	DownloadMedia(ctx context.Context, link string) (*models.Media, error)
}

type MessageService interface {
	SendMedia(ctx context.Context, chatID int64, media *models.Media) error
}

type Handler struct {
	downloaderService DownloaderService
	messageService    MessageService
}

func NewHandler(downloaderService DownloaderService, messageService MessageService) *Handler {
	return &Handler{
		downloaderService: downloaderService,
		messageService:    messageService,
	}
}

func (h *Handler) HandleUpdate(botAPI *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	ctx := context.Background()

	chatID := update.Message.Chat.ID
	text := update.Message.Text

	// Обработка команд
	if update.Message.IsCommand() {
		handleCommand(botAPI, update.Message)
		return
	}

	// Проверка наличия ссылки TikTok
	if sanitizer.IsTikTokLink(text) {
		// Загрузка медиа
		media, err := h.downloaderService.DownloadMedia(ctx, text)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка при загрузке контента.")
			log.Printf("Error downloading media: %v", err)
			_, _ = botAPI.Send(msg)
			return
		}
		// Отправка медиа
		err = h.messageService.SendMedia(ctx, chatID, media)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка при отправке контента.")
			log.Printf("Error downloading media: %v", err)
			_, _ = botAPI.Send(msg)
		}
	}
}
