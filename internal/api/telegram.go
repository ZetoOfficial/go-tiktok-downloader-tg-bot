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
	params, err := newVideoParams(chatID, options)
	if err != nil {
		return err
	}
	files := []tgbotapi.RequestFile{{
		Name: "video",
		Data: tgbotapi.FileBytes{Name: fileName, Bytes: data},
	}}
	_, err = ta.Bot.UploadFiles("sendVideo", params, files)
	return err
}

func (ta *TelegramAdapter) SendPhotoFile(chatID int64, fileName string, data []byte, opts ...models.SendOption) error {
	options := ta.applyOptions(opts)
	photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileBytes{Name: fileName, Bytes: data})
	if options.ReplyToMessageID > 0 {
		photoMsg.ReplyToMessageID = options.ReplyToMessageID
	}
	photoMsg.Caption, photoMsg.CaptionEntities = buildCaption(options)
	_, err := ta.Bot.Send(photoMsg)
	return err
}

// newVideoParams собирает параметры sendVideo вручную: VideoConfig из
// telegram-bot-api v5 не умеет передавать width и height, а без них
// Telegram показывает ролик квадратным.
func newVideoParams(chatID int64, options models.SendOptions) (tgbotapi.Params, error) {
	params := tgbotapi.Params{}
	params.AddNonZero64("chat_id", chatID)
	params.AddNonZero("reply_to_message_id", options.ReplyToMessageID)

	caption, entities := buildCaption(options)
	params.AddNonEmpty("caption", caption)
	if len(entities) > 0 {
		if err := params.AddInterface("caption_entities", entities); err != nil {
			return nil, err
		}
	}

	params.AddNonZero("width", options.Video.Width)
	params.AddNonZero("height", options.Video.Height)
	params.AddNonZero("duration", options.Video.Duration)
	params.AddBool("supports_streaming", true)
	return params, nil
}

func buildCaption(options models.SendOptions) (string, []tgbotapi.MessageEntity) {
	var caption string
	var entities []tgbotapi.MessageEntity
	if options.SourceName != "" && options.SourceURL != "" {
		caption = options.SourceName
		entities = []tgbotapi.MessageEntity{{
			Type:   "text_link",
			Offset: 0,
			Length: len(options.SourceName),
			URL:    options.SourceURL,
		}}
	}
	if options.Caption != "" {
		if caption != "" {
			caption += "\n"
		}
		caption += options.Caption
	}
	return caption, entities
}

func (ta *TelegramAdapter) SendMediaGroup(chatID int64, media []models.MediaInput, opts ...models.SendOption) error {
	if len(media) == 0 {
		return nil
	}

	options := ta.applyOptions(opts)

	var mediaGroup []interface{}
	for i, m := range media {
		file := tgbotapi.FileBytes{
			Name:  m.FileName,
			Bytes: m.Data,
		}
		caption := m.Caption
		var entities []tgbotapi.MessageEntity
		if i == 0 && options.SourceName != "" && options.SourceURL != "" {
			caption = options.SourceName
			if m.Caption != "" {
				caption += "\n" + m.Caption
			}
			entities = []tgbotapi.MessageEntity{{
				Type: "text_link", Offset: 0, Length: len(options.SourceName), URL: options.SourceURL,
			}}
		}
		switch m.Type {
		case "photo":
			inputMedia := tgbotapi.NewInputMediaPhoto(file)
			inputMedia.Caption = caption
			inputMedia.CaptionEntities = entities
			mediaGroup = append(mediaGroup, inputMedia)
		case "video":
			inputMedia := tgbotapi.NewInputMediaVideo(file)
			inputMedia.Caption = caption
			inputMedia.CaptionEntities = entities
			inputMedia.Width = m.Video.Width
			inputMedia.Height = m.Video.Height
			inputMedia.Duration = m.Video.Duration
			inputMedia.SupportsStreaming = true
			mediaGroup = append(mediaGroup, inputMedia)
		default:
			return errors.New("unsupported media type: " + m.Type)
		}
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
