package clients

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsInstagramPostLink(t *testing.T) {
	tests := []struct {
		link string
		want bool
	}{
		{"https://instagram.com/p/DdsrIIZMD7i/", true},
		{"https://www.instagram.com/p/DdsrIIZMD7i/?igsh=abc", true},
		{"https://m.instagram.com/p/DdsrIIZMD7i", true},
		{"https://instagram.com/reel/DdsrIIZMD7i/", false},
		{"https://instagram.com/p/DdsrIIZMD7i/extra", false},
		{"https://example.com/p/DdsrIIZMD7i/", false},
	}
	for _, tt := range tests {
		if got := isInstagramPostLink(tt.link); got != tt.want {
			t.Errorf("isInstagramPostLink(%q) = %t, want %t", tt.link, got, tt.want)
		}
	}
}

func TestCollectGalleryMediaKeepsCarouselOrder(t *testing.T) {
	dir := t.TempDir()
	for name, content := range map[string]string{
		"002.mp4": "second", "001.jpg": "first", "003.jpg": "third", "note.json": "ignore",
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	resp, err := collectGalleryMedia(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Items) != 3 {
		t.Fatalf("got %d items, want 3", len(resp.Items))
	}
	for i, want := range []struct{ name, kind, data string }{
		{"001.jpg", "photo", "first"},
		{"002.mp4", "video", "second"},
		{"003.jpg", "photo", "third"},
	} {
		got := resp.Items[i]
		if got.FileName != want.name || got.Type != want.kind || string(got.Data) != want.data {
			t.Errorf("item %d = %+v, want %+v", i, got, want)
		}
	}
}
