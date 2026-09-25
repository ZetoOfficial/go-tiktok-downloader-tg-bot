package clients

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name string, data []byte) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestCollectMedia_SingleVideo(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "clip.mp4", []byte("video-bytes"))

	resp, err := collectMedia(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.FileName != "clip.mp4" {
		t.Errorf("FileName = %q, want %q", resp.FileName, "clip.mp4")
	}
	if string(resp.Data) != "video-bytes" {
		t.Errorf("Data = %q, want %q", resp.Data, "video-bytes")
	}
	if len(resp.Photos) != 0 {
		t.Errorf("Photos = %v, want empty", resp.Photos)
	}
}

func TestCollectMedia_PhotoSlideshow(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "001.jpg", []byte("img-1"))
	writeFile(t, dir, "002.jpg", []byte("img-2"))
	writeFile(t, dir, "003.webp", []byte("img-3"))

	resp, err := collectMedia(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("Data = %q, want empty", resp.Data)
	}
	if len(resp.Photos) != 3 {
		t.Fatalf("len(Photos) = %d, want 3", len(resp.Photos))
	}
	if string(resp.Photos["001.jpg"]) != "img-1" {
		t.Errorf("Photos[001.jpg] = %q, want %q", resp.Photos["001.jpg"], "img-1")
	}
	if string(resp.Photos["003.webp"]) != "img-3" {
		t.Errorf("Photos[003.webp] = %q, want %q", resp.Photos["003.webp"], "img-3")
	}
}

func TestCollectMedia_IgnoresAudioInSlideshow(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "001.jpg", []byte("img-1"))
	writeFile(t, dir, "002.jpg", []byte("img-2"))
	writeFile(t, dir, "track.mp3", []byte("audio-bytes"))

	resp, err := collectMedia(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Photos) != 2 {
		t.Errorf("len(Photos) = %d, want 2 (audio must be ignored)", len(resp.Photos))
	}
}

func TestCollectMedia_VideoTakesPrecedenceOverImages(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "clip.mp4", []byte("video-bytes"))
	writeFile(t, dir, "clip.jpg", []byte("thumbnail")) // stray thumbnail

	resp, err := collectMedia(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(resp.Data) != "video-bytes" {
		t.Errorf("Data = %q, want video bytes", resp.Data)
	}
	if len(resp.Photos) != 0 {
		t.Errorf("Photos = %v, want empty when a video is present", resp.Photos)
	}
}

func TestCollectMedia_EmptyDirReturnsError(t *testing.T) {
	dir := t.TempDir()

	if _, err := collectMedia(dir); err == nil {
		t.Fatal("expected error for empty dir, got nil")
	}
}

func TestIsYouTubeLink(t *testing.T) {
	tests := []struct {
		link string
		want bool
	}{
		{"https://www.youtube.com/shorts/abc123", true},
		{"https://youtu.be/abc123", true},
		{"https://www.youtube.com/watch?v=abc123", false},
		{"https://tiktok.com/@user/video/youtube.com/shorts/abc123", false},
	}
	for _, tt := range tests {
		if got := isYouTubeLink(tt.link); got != tt.want {
			t.Errorf("isYouTubeLink(%q) = %t, want %t", tt.link, got, tt.want)
		}
	}
}
