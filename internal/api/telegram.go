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

func (ta *TelegramAdapter) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := ta.Bot.Send(msg)
	return err
}

func (ta *TelegramAdapter) SendVideoFile(chatID int64, fileName string, data []byte) error {
	videoFile := tgbotapi.FileBytes{
		Name:  fileName,
		Bytes: data,
	}
	videoMsg := tgbotapi.NewVideo(chatID, videoFile)
	_, err := ta.Bot.Send(videoMsg)
	return err
}

func (ta *TelegramAdapter) SendMediaGroup(chatID int64, media []models.MediaInput) error {
	if len(media) == 0 {
		return nil
	}

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
		ChatID: chatID,
		Media:  mediaGroup,
	}

	_, err := ta.Bot.SendMediaGroup(config)
	return err
}
