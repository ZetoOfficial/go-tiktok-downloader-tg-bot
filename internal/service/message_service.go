package service

import (
	"context"
	"fmt"
	"sort"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
)

type TelegramSender interface {
	SendMessage(chatID int64, text string, opts ...models.SendOption) error
	SendVideoFile(chatID int64, fileName string, data []byte, opts ...models.SendOption) error
	SendPhotoFile(chatID int64, fileName string, data []byte, opts ...models.SendOption) error
	SendMediaGroup(chatID int64, media []models.MediaInput, opts ...models.SendOption) error
}

type MessageService struct {
	sender TelegramSender
}

func NewMessageService(sender TelegramSender) *MessageService {
	return &MessageService{sender: sender}
}

func (ms *MessageService) SendMedia(ctx context.Context, chatID int64, media *models.Media, opts ...models.SendOption) error {
	if len(media.Items) > 0 {
		return ms.sendItems(ctx, chatID, media.Items, opts...)
	}
	if len(media.VideoData) > 0 {
		return ms.sender.SendVideoFile(chatID, media.VideoName, media.VideoData, opts...)
	}

	// Отправка изображений, если они есть
	if len(media.Photos) > 0 {
		// Сортировка имен файлов для последовательности
		var keys []string
		for k := range media.Photos {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		items := make([]models.MediaInput, 0, len(keys))
		for _, key := range keys {
			items = append(items, models.MediaInput{Type: "photo", FileName: key, Data: media.Photos[key]})
		}
		return ms.sendItems(ctx, chatID, items, opts...)
	}
	return nil
}

// sendItems отправляет медиа в исходном порядке, группами по 10.
func (ms *MessageService) sendItems(ctx context.Context, chatID int64, items []models.MediaInput, opts ...models.SendOption) error {
	totalParts := (len(items) + 9) / 10
	for i := 0; i < len(items); i += 10 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		end := i + 10
		if end > len(items) {
			end = len(items)
		}
		inputs := append([]models.MediaInput(nil), items[i:end]...)
		partNumber := i/10 + 1
		if totalParts > 1 {
			inputs[0].Caption = fmt.Sprintf("Часть %d/%d", partNumber, totalParts)
		}
		if len(inputs) == 1 {
			item := inputs[0]
			singleOpts := append(append([]models.SendOption(nil), opts...), models.WithCaption(item.Caption))
			switch item.Type {
			case "photo":
				if err := ms.sender.SendPhotoFile(chatID, item.FileName, item.Data, singleOpts...); err != nil {
					return err
				}
			case "video":
				if err := ms.sender.SendVideoFile(chatID, item.FileName, item.Data, singleOpts...); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported media type: %s", item.Type)
			}
			continue
		}
		if err := ms.sender.SendMediaGroup(chatID, inputs, opts...); err != nil {
			return err
		}
	}
	return nil
}
