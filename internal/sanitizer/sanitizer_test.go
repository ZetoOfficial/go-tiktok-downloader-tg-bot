package sanitizer

import "testing"

func TestIsInstagramReelLink(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{name: "www reel", text: "https://www.instagram.com/reel/C1Ab_2-x/", want: true},
		{name: "mobile reel", text: "https://m.instagram.com/reel/C1Ab_2-x/", want: true},
		{name: "query parameters", text: "https://instagram.com/reel/C1Ab_2-x?igsh=YWJjZA==", want: true},
		{name: "link in message", text: "Watch this https://www.instagram.com/reel/C1Ab_2-x/ now", want: true},
		{name: "profile", text: "https://www.instagram.com/example/", want: false},
		{name: "post", text: "https://www.instagram.com/p/C1Ab_2-x/", want: false},
		{name: "story", text: "https://www.instagram.com/stories/example/123/", want: false},
		{name: "reels page", text: "https://www.instagram.com/reels/", want: false},
		{name: "different domain", text: "https://example.com/reel/C1Ab_2-x/", want: false},
		{name: "no link", text: "send me an Instagram Reel", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInstagramReelLink(tt.text); got != tt.want {
				t.Errorf("IsInstagramReelLink(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestExistingLinkSanitizers(t *testing.T) {
	if !IsTikTokLink("https://www.tiktok.com/@example/video/123") {
		t.Error("TikTok link must remain supported")
	}
	if !IsYouTubeShortsLink("https://www.youtube.com/shorts/example") {
		t.Error("YouTube Shorts link must remain supported")
	}
}
