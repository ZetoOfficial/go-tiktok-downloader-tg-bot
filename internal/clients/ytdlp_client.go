package clients

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
)

// ytdlpTimeout ограничивает время работы одного вызова yt-dlp.
const ytdlpTimeout = 3 * time.Minute

// maxFileSizeArg задаёт лимит размера скачиваемого файла для yt-dlp.
// Telegram не принимает от ботов файлы больше ~50 МБ.
const maxFileSizeArg = "50M"

var videoExtensions = map[string]bool{
	".mp4": true, ".webm": true, ".mkv": true, ".mov": true,
}

var imageExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true,
}

// YtDlpClient скачивает видео и фото-слайдшоу (TikTok, YouTube Shorts и др.),
// вызывая бинарь yt-dlp. Реализует интерфейс service.MediaDownloader.
type YtDlpClient struct {
	binary string
}

func NewYtDlpClient() *YtDlpClient {
	return &YtDlpClient{binary: "yt-dlp"}
}

func (c *YtDlpClient) Download(ctx context.Context, link string, _ ...models.DownloadOption) (*models.DownloadResponse, error) {
	if link == "" {
		return nil, fmt.Errorf("empty link")
	}

	tmpDir, err := os.MkdirTemp("", "ytdlp-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	runCtx, cancel := context.WithTimeout(ctx, ytdlpTimeout)
	defer cancel()

	outputTemplate := filepath.Join(tmpDir, "%(autonumber)03d-%(title).50s.%(ext)s")
	args := []string{
		"--no-progress",
		"--js-runtimes", "deno",
		"--max-filesize", maxFileSizeArg,
		"-S", "res:720,ext:mp4:m4a,filesize:50M",
		"--merge-output-format", "mp4",
		"--socket-timeout", "30",
		"-o", outputTemplate,
		"--", link,
	}

	var stderr bytes.Buffer
	cmd := exec.CommandContext(runCtx, c.binary, args...)
	cmd.Stderr = &stderr

	log.Printf("yt-dlp download: %s", link)
	if err := cmd.Run(); err != nil {
		if runCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("yt-dlp timed out after %s", ytdlpTimeout)
		}
		return nil, fmt.Errorf("yt-dlp failed: %w: %s", err, tail(stderr.String(), 500))
	}
	if warnings := strings.TrimSpace(stderr.String()); warnings != "" {
		log.Printf("yt-dlp warnings: %s", warnings)
	}

	return collectMedia(tmpDir)
}

// collectMedia читает скачанные yt-dlp файлы из dir и раскладывает их в
// DownloadResponse: если есть видео — возвращает видео (картинки-превью
// игнорируются); иначе собирает картинки в альбом. Аудио и прочие файлы
// (json, .description) игнорируются.
func collectMedia(dir string) (*models.DownloadResponse, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read output dir: %w", err)
	}

	var videoNames []string
	var imageNames []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		switch {
		case videoExtensions[ext]:
			videoNames = append(videoNames, e.Name())
		case imageExtensions[ext]:
			imageNames = append(imageNames, e.Name())
		}
	}

	// Видео имеет приоритет: если yt-dlp скачал ролик, возможные превью-картинки
	// рядом с ним не считаем контентом.
	if len(videoNames) > 0 {
		sort.Strings(videoNames)
		name := videoNames[0]
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("read video %s: %w", name, err)
		}
		return &models.DownloadResponse{FileName: name, Data: data}, nil
	}

	if len(imageNames) > 0 {
		photos := make(map[string][]byte, len(imageNames))
		for _, name := range imageNames {
			data, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				return nil, fmt.Errorf("read image %s: %w", name, err)
			}
			photos[name] = data
		}
		return &models.DownloadResponse{Photos: photos}, nil
	}

	return nil, fmt.Errorf("no media downloaded")
}

// tail возвращает последние n символов строки (для компактного вывода ошибок).
func tail(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return "..." + s[len(s)-n:]
}
