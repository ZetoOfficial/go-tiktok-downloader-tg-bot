package api

import (
	"testing"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
)

func TestNewVideoConfigSourceLink(t *testing.T) {
	config := newVideoConfig(123, "clip.mp4", []byte("video"), models.SendOptions{
		ReplyToMessageID: 42,
		SourceName:       "YouTube",
		SourceURL:        "https://www.youtube.com/shorts/abc123?si=a&b=c",
	})
	if config.Caption != "YouTube" || config.ReplyToMessageID != 42 {
		t.Fatalf("video caption = %q, reply = %d", config.Caption, config.ReplyToMessageID)
	}
	if len(config.CaptionEntities) != 1 {
		t.Fatalf("caption entities = %v, want one link", config.CaptionEntities)
	}
	entity := config.CaptionEntities[0]
	if entity.Type != "text_link" || entity.Offset != 0 || entity.Length != len("YouTube") || entity.URL != "https://www.youtube.com/shorts/abc123?si=a&b=c" {
		t.Errorf("caption link = %+v", entity)
	}
}
