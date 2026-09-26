package service

import (
	"context"
	"testing"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
)

type recordingSender struct {
	photos []string
	videos []string
	groups [][]models.MediaInput
}

func (*recordingSender) SendMessage(int64, string, ...models.SendOption) error { return nil }
func (s *recordingSender) SendPhotoFile(_ int64, name string, _ []byte, _ ...models.SendOption) error {
	s.photos = append(s.photos, name)
	return nil
}
func (s *recordingSender) SendVideoFile(_ int64, name string, _ []byte, _ ...models.SendOption) error {
	s.videos = append(s.videos, name)
	return nil
}
func (s *recordingSender) SendMediaGroup(_ int64, items []models.MediaInput, _ ...models.SendOption) error {
	s.groups = append(s.groups, items)
	return nil
}

func TestSendMediaSinglePhoto(t *testing.T) {
	sender := &recordingSender{}
	media := &models.Media{Items: []models.MediaInput{{Type: "photo", FileName: "001.jpg", Data: []byte("photo")}}}
	if err := NewMessageService(sender).SendMedia(context.Background(), 1, media); err != nil {
		t.Fatal(err)
	}
	if len(sender.photos) != 1 || sender.photos[0] != "001.jpg" || len(sender.groups) != 0 {
		t.Fatalf("photos = %v, groups = %v", sender.photos, sender.groups)
	}
}

func TestSendMediaMixedCarouselWithRemainder(t *testing.T) {
	sender := &recordingSender{}
	items := make([]models.MediaInput, 11)
	for i := range items {
		items[i] = models.MediaInput{Type: "photo", FileName: "photo.jpg"}
	}
	items[1] = models.MediaInput{Type: "video", FileName: "video.mp4"}
	items[10] = models.MediaInput{Type: "video", FileName: "last.mp4"}
	if err := NewMessageService(sender).SendMedia(context.Background(), 1, &models.Media{Items: items}); err != nil {
		t.Fatal(err)
	}
	if len(sender.groups) != 1 || len(sender.groups[0]) != 10 || sender.groups[0][1].Type != "video" {
		t.Fatalf("groups = %v", sender.groups)
	}
	if len(sender.videos) != 1 || sender.videos[0] != "last.mp4" {
		t.Fatalf("videos = %v", sender.videos)
	}
}
