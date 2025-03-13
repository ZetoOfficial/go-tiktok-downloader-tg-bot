package clients

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/config"
	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
	"github.com/go-resty/resty/v2"
)

type DouyinClient struct {
	client *resty.Client
	apiURL *url.URL
}

func NewDouyinClient(cfg *config.Config) (*DouyinClient, error) {
	parsedURL, err := url.Parse(cfg.DouyinAPI)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	client := resty.New().
		SetBaseURL(parsedURL.String()).
		SetTimeout(30 * time.Second).
		SetRetryCount(3)

	return &DouyinClient{
		client: client,
		apiURL: parsedURL,
	}, nil
}

// Download осуществляет запрос к эндпоинту /download с дополнительными параметрами.
// Если ответ представляет zip-архив, он распаковывается и извлекаются фотографии.
func (c *DouyinClient) Download(ctx context.Context, link string, options ...models.DownloadOption) (*models.DownloadResponse, error) {
	if link == "" {
		return nil, fmt.Errorf("empty link")
	}

	params := map[string]string{
		"url":            link,
		"prefix":         "false",
		"with_watermark": "true",
	}

	for _, opt := range options {
		params[opt.Key] = fmt.Sprintf("%t", opt.Value)
	}

	startTime := time.Now()

	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		Get("/download")

	duration := time.Since(startTime) // Вычисляем длительность запроса
	log.Printf("Download request completed in %v", duration)

	if err != nil {
		return nil, fmt.Errorf("download: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("api: %s", resp.Status())
	}

	// Попытка извлечь имя файла из заголовка Content-Disposition.
	filename := extractFileName(resp.Header().Get("Content-Disposition"))
	if filename == "" {
		filename = "downloaded_file"
	}

	// Если сервер возвращает zip-архив (определяем по заголовку или расширению файла)
	contentType := resp.Header().Get("Content-Type")
	if contentType == "application/zip" || strings.HasSuffix(filename, ".zip") {
		photos, err := extractPhotosFromZip(resp.Body())
		if err != nil {
			return nil, fmt.Errorf("extract zip: %w", err)
		}
		return &models.DownloadResponse{
			FileName: filename,
			Photos:   photos,
		}, nil
	}

	// Иначе возвращаем обычный файл (например, видео)
	return &models.DownloadResponse{
		FileName: filename,
		Data:     resp.Body(),
	}, nil
}

// extractFileName пытается извлечь имя файла из заголовка Content-Disposition.
func extractFileName(cd string) string {
	if cd == "" {
		return ""
	}
	// Ожидаемый формат: attachment; filename="file.mp4"
	parts := strings.Split(cd, "filename=")
	if len(parts) < 2 {
		return ""
	}
	return strings.Trim(parts[1], `"`)
}

// extractPhotosFromZip распаковывает zip-архив и возвращает все файлы с изображениями.
func extractPhotosFromZip(data []byte) (map[string][]byte, error) {
	photos := make(map[string][]byte)
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	for _, file := range reader.File {
		// Пропускаем директории
		if file.FileInfo().IsDir() {
			continue
		}
		// Обрабатываем только файлы с расширениями изображений
		lowerName := strings.ToLower(file.Name)
		if strings.HasSuffix(lowerName, ".jpg") ||
			strings.HasSuffix(lowerName, ".jpeg") ||
			strings.HasSuffix(lowerName, ".png") ||
			strings.HasSuffix(lowerName, ".webp") ||
			strings.HasSuffix(lowerName, ".gif") {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, err
			}
			photos[file.Name] = content
		}
	}
	return photos, nil
}
