package api

import (
	"encoding/json"
	"testing"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestNewVideoParamsSourceLink(t *testing.T) {
	params, err := newVideoParams(123, models.SendOptions{
		ReplyToMessageID: 42,
		SourceName:       "YouTube",
		SourceURL:        "https://www.youtube.com/shorts/abc123?si=a&b=c",
	})
	if err != nil {
		t.Fatal(err)
	}
	if params["chat_id"] != "123" || params["caption"] != "YouTube" || params["reply_to_message_id"] != "42" {
		t.Fatalf("params = %v", params)
	}
	var entities []tgbotapi.MessageEntity
	if err := json.Unmarshal([]byte(params["caption_entities"]), &entities); err != nil {
		t.Fatalf("caption_entities = %q: %v", params["caption_entities"], err)
	}
	if len(entities) != 1 {
		t.Fatalf("caption entities = %v, want one link", entities)
	}
	entity := entities[0]
	if entity.Type != "text_link" || entity.Offset != 0 || entity.Length != len("YouTube") || entity.URL != "https://www.youtube.com/shorts/abc123?si=a&b=c" {
		t.Errorf("caption link = %+v", entity)
	}
	for _, key := range []string{"width", "height", "duration"} {
		if _, ok := params[key]; ok {
			t.Errorf("%s set without video meta: %v", key, params)
		}
	}
}

func TestNewVideoParamsVideoMeta(t *testing.T) {
	params, err := newVideoParams(1, models.SendOptions{
		Caption: "Часть 1/2",
		Video:   models.VideoMeta{Width: 720, Height: 1280, Duration: 15},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{
		"width": "720", "height": "1280", "duration": "15",
		"supports_streaming": "true", "caption": "Часть 1/2",
	}
	for key, value := range want {
		if params[key] != value {
			t.Errorf("%s = %q, want %q", key, params[key], value)
		}
	}
	if _, ok := params["caption_entities"]; ok {
		t.Errorf("caption_entities set without source: %v", params)
	}
}
