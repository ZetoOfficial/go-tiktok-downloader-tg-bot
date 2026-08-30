package sanitizer

import "regexp"

var (
	tikTokPattern        = regexp.MustCompile(`(https?://)?(www\.)?(tiktok\.com|douyin\.com)/\S+`)
	youTubeShortsPattern = regexp.MustCompile(`(https?://)?(www\.|m\.)?(youtube\.com/shorts/|youtu\.be/)\S+`)
	instagramReelPattern = regexp.MustCompile(`(?i)https?://(?:www\.|m\.)?instagram\.com/reel/[a-z0-9_-]+/?(?:\?[^\s]*)?(?:\s|$)`)
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
