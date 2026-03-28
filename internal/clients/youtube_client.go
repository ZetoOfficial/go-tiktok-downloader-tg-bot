package clients

import (
	"context"
	"fmt"
	"io"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
	"github.com/kkdai/youtube/v2"
)

const maxVideoSize = 50 * 1024 * 1024 // 50 MB

type YouTubeClient struct {
	client youtube.Client
}

func NewYouTubeClient() *YouTubeClient {
	return &YouTubeClient{
		client: youtube.Client{},
	}
}

func (c *YouTubeClient) Download(ctx context.Context, link string, options ...models.DownloadOption) (*models.DownloadResponse, error) {
	video, err := c.client.GetVideoContext(ctx, link)
	if err != nil {
		return nil, fmt.Errorf("get video info: %w", err)
	}

	format := selectFormat(video.Formats)
	if format == nil {
		return nil, fmt.Errorf("no suitable video format found (mp4 with audio, under 50MB)")
	}

	stream, _, err := c.client.GetStreamContext(ctx, video, format)
	if err != nil {
		return nil, fmt.Errorf("get stream: %w", err)
	}
	defer stream.Close()

	data, err := io.ReadAll(io.LimitReader(stream, maxVideoSize+1))
	if err != nil {
		return nil, fmt.Errorf("read stream: %w", err)
	}
	if len(data) > maxVideoSize {
		return nil, fmt.Errorf("video is too large (>50MB)")
	}

	filename := video.Title + ".mp4"

	return &models.DownloadResponse{
		FileName: filename,
		Data:     data,
	}, nil
}

func selectFormat(formats youtube.FormatList) *youtube.Format {
	// Prefer mp4 formats with audio
	filtered := formats.Type("video/mp4").WithAudioChannels()
	if len(filtered) == 0 {
		filtered = formats.WithAudioChannels()
	}
	if len(filtered) == 0 {
		return nil
	}

	// Sort by quality descending, pick the best that fits under 50MB
	filtered.Sort()
	for i := len(filtered) - 1; i >= 0; i-- {
		if filtered[i].ContentLength > 0 && filtered[i].ContentLength <= int64(maxVideoSize) {
			return &filtered[i]
		}
	}

	// If ContentLength is unknown, try the smallest available
	return &filtered[0]
}
