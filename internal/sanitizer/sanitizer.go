package sanitizer

import "regexp"

func IsTikTokLink(text string) bool {
	pattern := `(https?://)?(www\.)?(tiktok\.com|douyin\.com)/\S+`
	re := regexp.MustCompile(pattern)
	return re.MatchString(text)
}
