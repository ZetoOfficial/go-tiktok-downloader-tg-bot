package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/sanitizer"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DownloaderService interface {
	DownloadMedia(ctx context.Context, link string) (*models.Media, error)
}

type MessageService interface {
	SendMedia(ctx context.Context, chatID int64, media *models.Media, opts ...models.SendOption) error
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

func (h *Handler) HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	ctx := context.Background()

	msg := update.Message
	chatID := msg.Chat.ID
	messageID := msg.MessageID
	text := msg.Text

	if update.Message.IsCommand() {
		handleCommand(bot, update.Message)
		return
	}

	sourceName, sourceURL, supported := sanitizer.SourceLink(text)
	isInstagram := sourceName == "Instagram"

	// Проверка наличия поддерживаемой ссылки
	if supported {
		// 👀 реакция "смотрю"
		h.setReaction(bot, chatID, messageID, "👀")

		media, err := h.downloaderService.DownloadMedia(ctx, sourceURL)
		if err != nil {
			userMessage := "Ошибка при загрузке контента."
			if isInstagram {
				userMessage = "Не удалось скачать публикацию Instagram. Поддерживаются только публичные посты и Reels: Instagram мог потребовать вход, ограничить доступ или публикация недоступна."
			}
			h.replyWithError(bot, chatID, messageID, userMessage, err)
			// 👎 реакция "плаки-плаки"
			h.setReaction(bot, chatID, messageID, "👎")
			return
		}

		if err := h.messageService.SendMedia(ctx, chatID, media, models.WithReplyTo(messageID), models.WithSource(sourceName, sourceURL)); err != nil {
			h.replyWithError(bot, chatID, messageID, "Ошибка при отправке контента.", err)
			// 👎 реакция "плаки-плаки"
			h.setReaction(bot, chatID, messageID, "👎")
			return
		}

		// 👍 реакция "всё ок"
		h.setReaction(bot, chatID, messageID, "👍")
	}
}

func (h *Handler) replyWithError(bot *tgbotapi.BotAPI, chatID int64, replyTo int, userMsg string, err error) {
	log.Printf("Handler error: %v", err)

	msg := tgbotapi.NewMessage(chatID, userMsg)
	msg.ReplyToMessageID = replyTo
	_, _ = bot.Send(msg)
}

func (h *Handler) setReaction(bot *tgbotapi.BotAPI, chatID int64, messageID int, emoji string) {
	reaction := []models.EmojiReactionPayload{
		{Type: "emoji", Emoji: emoji},
	}

	reactionJSON, err := json.Marshal(reaction)
	if err != nil {
		log.Printf("setReaction: marshal error: %v", err)
		return
	}

	params := tgbotapi.Params{
		"chat_id":    fmt.Sprintf("%d", chatID),
		"message_id": fmt.Sprintf("%d", messageID),
		"reaction":   string(reactionJSON),
	}

	if _, err := bot.MakeRequest("setMessageReaction", params); err != nil {
		log.Printf("setReaction: request error: %v", err)
	}
}
