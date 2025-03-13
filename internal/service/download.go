package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
)

const maxVideoSize = 50 * 1024 * 1024 // 50 MB

type DouyinClient interface {
	Download(ctx context.Context, link string, options ...models.DownloadOption) (*models.DownloadResponse, error)
}

type DownloadService struct {
	douyinClient DouyinClient
}

func NewDownloadService(douyinClient DouyinClient) *DownloadService {
	return &DownloadService{
		douyinClient: douyinClient,
	}
}

func (d *DownloadService) DownloadMedia(ctx context.Context, link string) (*models.Media, error) {
	log.Print("Downloading media...")
	resp, err := d.douyinClient.Download(ctx, link)
	if err != nil {
		return nil, fmt.Errorf("download media: %w", err)
	}

	if len(resp.Data) > 0 {
		if len(resp.Data) > maxVideoSize {
			return nil, errors.New("video is too large")
		}
		return &models.Media{
			VideoData: resp.Data,
			VideoName: resp.FileName,
		}, nil
	}

	if len(resp.Photos) > 0 {
		return &models.Media{
			Photos: resp.Photos,
		}, nil
	}

	return nil, errors.New("media not found")
}
