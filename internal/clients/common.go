package clients

import (
	"archive/zip"
	"bytes"
	"io"
	"strconv"
	"strings"
)

func extractFileName(cd string) string {
	if cd == "" {
		return ""
	}
	parts := strings.Split(cd, "filename=")
	if len(parts) < 2 {
		return ""
	}
	return strings.Trim(parts[1], `"`)
}

func extractPhotosFromZip(data []byte) (map[string][]byte, error) {
	photos := make(map[string][]byte)
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}
		lowerName := strings.ToLower(file.Name)
		if strings.HasSuffix(lowerName, ".jpg") ||
			strings.HasSuffix(lowerName, ".jpeg") ||
			strings.HasSuffix(lowerName, ".png") ||
			strings.HasSuffix(lowerName, ".webp") ||
			strings.HasSuffix(lowerName, ".gif") {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, err
			}
			photos[file.Name] = content
		}
	}
	return photos, nil
}

func parseContentLength(cl string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(cl))
}
