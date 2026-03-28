package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/sanitizer"
)

const maxVideoSize = 50 * 1024 * 1024 // 50 MB

type MediaDownloader interface {
	Download(ctx context.Context, link string, options ...models.DownloadOption) (*models.DownloadResponse, error)
}

type DownloadService struct {
	douyinClient  MediaDownloader
	youtubeClient MediaDownloader
}

func NewDownloadService(douyinClient, youtubeClient MediaDownloader) *DownloadService {
	return &DownloadService{
		douyinClient:  douyinClient,
		youtubeClient: youtubeClient,
	}
}

func (d *DownloadService) DownloadMedia(ctx context.Context, link string) (*models.Media, error) {
	log.Print("Downloading media...")

	var resp *models.DownloadResponse
	var err error

	switch {
	case sanitizer.IsYouTubeShortsLink(link):
		resp, err = d.youtubeClient.Download(ctx, link)
	default:
		resp, err = d.douyinClient.Download(ctx, link)
	}
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
