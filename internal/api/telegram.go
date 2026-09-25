package api

import (
	"errors"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramAdapter struct {
	Bot *tgbotapi.BotAPI
}

func NewTelegramAdapter(bot *tgbotapi.BotAPI) *TelegramAdapter {
	return &TelegramAdapter{Bot: bot}
}

func (ta *TelegramAdapter) SendMessage(chatID int64, text string, opts ...models.SendOption) error {
	options := ta.applyOptions(opts)

	msg := tgbotapi.NewMessage(chatID, text)
	if options.ReplyToMessageID > 0 {
		msg.ReplyToMessageID = options.ReplyToMessageID
	}

	_, err := ta.Bot.Send(msg)
	return err
}

func (ta *TelegramAdapter) SendVideoFile(chatID int64, fileName string, data []byte, opts ...models.SendOption) error {
	options := ta.applyOptions(opts)
	videoMsg := newVideoConfig(chatID, fileName, data, options)
	_, err := ta.Bot.Send(videoMsg)
	return err
}

func newVideoConfig(chatID int64, fileName string, data []byte, options models.SendOptions) tgbotapi.VideoConfig {
	videoFile := tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: data,
	}
	videoMsg := tgbotapi.NewVideo(chatID, videoFile)
	if options.SourceName != "" && options.SourceURL != "" {
		videoMsg.Caption = options.SourceName
		videoMsg.CaptionEntities = []tgbotapi.MessageEntity{{
			Type:   "text_link",
			Offset: 0,
			Length: len(options.SourceName),
			URL:    options.SourceURL,
		}}
	}

	if options.ReplyToMessageID > 0 {
		videoMsg.ReplyToMessageID = options.ReplyToMessageID
	}

	return videoMsg
}

func (ta *TelegramAdapter) SendMediaGroup(chatID int64, media []models.MediaInput, opts ...models.SendOption) error {
	if len(media) == 0 {
		return nil
	}

	options := ta.applyOptions(opts)

	var mediaGroup []interface{}
	for _, m := range media {
		if m.Type != "photo" {
			continue
		}

		file := tgbotapi.FileBytes{
			Name:  m.FileName,
			Bytes: m.Data,
		}

		inputMedia := tgbotapi.NewInputMediaPhoto(file)
		inputMedia.Caption = m.Caption
		mediaGroup = append(mediaGroup, inputMedia)
	}

	if len(mediaGroup) == 0 {
		return errors.New("media group is empty")
	}

	config := tgbotapi.MediaGroupConfig{
		ChatID:           chatID,
		Media:            mediaGroup,
		ReplyToMessageID: options.ReplyToMessageID,
	}

	_, err := ta.Bot.SendMediaGroup(config)
	return err
}

func (ta *TelegramAdapter) applyOptions(opts []models.SendOption) models.SendOptions {
	options := models.SendOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
