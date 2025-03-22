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
	SendMediaGroup(chatID int64, media []models.MediaInput, opts ...models.SendOption) error
}

type MessageService struct {
	sender TelegramSender
}

func NewMessageService(sender TelegramSender) *MessageService {
	return &MessageService{sender: sender}
}

func (ms *MessageService) SendMedia(ctx context.Context, chatID int64, media *models.Media, opts ...models.SendOption) error {
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

		// Вычисление общего количества частей
		totalParts := (len(keys) + 9) / 10

		// Отправка изображений группами по 10
		for i := 0; i < len(keys); i += 10 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			var inputs []models.MediaInput
			groupSize := 10
			if i+groupSize > len(keys) {
				groupSize = len(keys) - i
			}
			partNumber := i/10 + 1
			for j := 0; j < groupSize; j++ {
				key := keys[i+j]
				input := models.MediaInput{
					Type:     "photo",
					FileName: key,
					Data:     media.Photos[key],
				}
				if j == 0 {
					input.Caption = fmt.Sprintf("Часть %d/%d", partNumber, totalParts)
				}
				inputs = append(inputs, input)
			}
			if err := ms.sender.SendMediaGroup(chatID, inputs, opts...); err != nil {
				return err
			}
		}
	}
	return nil
}
