package sanitizer

import (
	"regexp"
	"strings"
)

var (
	tikTokPattern        = regexp.MustCompile(`(https?://)?(www\.|m\.|vt\.|vm\.)?(tiktok\.com|douyin\.com)/\S+`)
	youTubeShortsPattern = regexp.MustCompile(`(https?://)?(www\.|m\.)?(youtube\.com/shorts/|youtu\.be/)\S+`)
	instagramReelPattern = regexp.MustCompile(`(?i)https?://(?:www\.|m\.)?instagram\.com/reel/[a-z0-9_-]+/?(?:\?[^\s]*)?(?:\s|$)`)
	instagramPostPattern = regexp.MustCompile(`(?i)https?://(?:www\.|m\.)?instagram\.com/p/[a-z0-9_-]+/?(?:\?[^\s]*)?(?:\s|$)`)
)

func IsTikTokLink(text string) bool {
	return tikTokPattern.MatchString(text)
}

func IsYouTubeShortsLink(text string) bool {
	return youTubeShortsPattern.MatchString(text)
}

func IsInstagramReelLink(text string) bool {
	return instagramReelPattern.MatchString(text)
}

// SourceLink возвращает платформу и поддерживаемую ссылку из сообщения.
func SourceLink(text string) (name, link string, ok bool) {
	for _, source := range []struct {
		name    string
		pattern *regexp.Regexp
	}{
		{"TikTok", tikTokPattern},
		{"YouTube", youTubeShortsPattern},
		{"Instagram", instagramReelPattern},
		{"Instagram", instagramPostPattern},
	} {
		link = strings.TrimRight(strings.TrimSpace(source.pattern.FindString(text)), ".,!?)")
		if link == "" {
			continue
		}
		if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") {
			link = "https://" + link
		}
		return source.name, link, true
	}
	return "", "", false
}
