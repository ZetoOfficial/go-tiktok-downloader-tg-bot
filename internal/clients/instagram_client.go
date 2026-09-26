package clients

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
)

func isInstagramPostLink(link string) bool {
	u, err := url.Parse(link)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return false
	}
	switch strings.ToLower(u.Hostname()) {
	case "instagram.com", "www.instagram.com", "m.instagram.com":
	default:
		return false
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	return len(parts) == 2 && parts[0] == "p" && parts[1] != ""
}

func (c *YtDlpClient) downloadInstagramPost(ctx context.Context, link string) (*models.DownloadResponse, error) {
	tmpDir, err := os.MkdirTemp("", "instagram-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	runCtx, cancel := context.WithTimeout(ctx, ytdlpTimeout)
	defer cancel()

	cmd := exec.CommandContext(runCtx, c.galleryBinary,
		"--config-ignore", "--filesize-max", maxFileSizeArg,
		"-D", tmpDir, "-f", "{num:03}.{extension}",
		"--", link,
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if runCtx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("gallery-dl timed out after %s", ytdlpTimeout)
		}
		return nil, fmt.Errorf("gallery-dl failed: %w: %s", err, tail(stderr.String(), 500))
	}
	return collectGalleryMedia(tmpDir)
}

func collectGalleryMedia(dir string) (*models.DownloadResponse, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read output dir: %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	items := make([]models.MediaInput, 0, len(names))
	for _, name := range names {
		ext := strings.ToLower(filepath.Ext(name))
		mediaType := ""
		switch {
		case videoExtensions[ext]:
			mediaType = "video"
		case imageExtensions[ext]:
			mediaType = "photo"
		default:
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", name, err)
		}
		items = append(items, models.MediaInput{Type: mediaType, FileName: name, Data: data})
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no media downloaded")
	}
	return &models.DownloadResponse{Items: items}, nil
}
