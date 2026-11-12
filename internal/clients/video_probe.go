package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
)

// probeTimeout ограничивает один вызов ffprobe.
const probeTimeout = 15 * time.Second

type ffprobeOutput struct {
	Streams []struct {
		Width             int               `json:"width"`
		Height            int               `json:"height"`
		Duration          string            `json:"duration"`
		SampleAspectRatio string            `json:"sample_aspect_ratio"`
		Tags              map[string]string `json:"tags"`
		SideDataList      []struct {
			Rotation float64 `json:"rotation"`
		} `json:"side_data_list"`
	} `json:"streams"`
	Format struct {
		Duration string `json:"duration"`
	} `json:"format"`
}

// probeVideo читает через ffprobe размеры видео так, как их увидит зритель:
// с учётом поворота и неквадратных пикселей.
func probeVideo(ctx context.Context, binary, path string) (models.VideoMeta, error) {
	runCtx, cancel := context.WithTimeout(ctx, probeTimeout)
	defer cancel()

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(runCtx, binary,
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,duration,sample_aspect_ratio:stream_tags=rotate:stream_side_data=rotation:format=duration",
		"-of", "json",
		path, // абсолютный путь во временной папке, с "-" не начинается
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return models.VideoMeta{}, fmt.Errorf("ffprobe failed: %w: %s", err, tail(stderr.String(), 300))
	}
	return parseFFprobe(stdout.Bytes())
}

func parseFFprobe(data []byte) (models.VideoMeta, error) {
	var out ffprobeOutput
	if err := json.Unmarshal(data, &out); err != nil {
		return models.VideoMeta{}, fmt.Errorf("parse ffprobe output: %w", err)
	}
	if len(out.Streams) == 0 {
		return models.VideoMeta{}, fmt.Errorf("no video stream")
	}
	s := out.Streams[0]
	width, height := s.Width, s.Height
	if width <= 0 || height <= 0 {
		return models.VideoMeta{}, fmt.Errorf("invalid video size %dx%d", width, height)
	}

	if num, den, ok := parseRatio(s.SampleAspectRatio); ok && num != den {
		width = int(math.Round(float64(width) * float64(num) / float64(den)))
	}

	rotation := 0.0
	for _, sd := range s.SideDataList {
		if sd.Rotation != 0 {
			rotation = sd.Rotation
			break
		}
	}
	if rotation == 0 {
		if r, err := strconv.ParseFloat(s.Tags["rotate"], 64); err == nil {
			rotation = r
		}
	}
	if int(math.Abs(math.Round(rotation)))%180 == 90 {
		width, height = height, width
	}

	duration := parseSeconds(s.Duration)
	if duration == 0 {
		duration = parseSeconds(out.Format.Duration)
	}
	return models.VideoMeta{Width: width, Height: height, Duration: duration}, nil
}

func parseRatio(s string) (int, int, bool) {
	a, b, ok := strings.Cut(s, ":")
	if !ok {
		return 0, 0, false
	}
	num, err1 := strconv.Atoi(a)
	den, err2 := strconv.Atoi(b)
	if err1 != nil || err2 != nil || num <= 0 || den <= 0 {
		return 0, 0, false
	}
	return num, den, true
}

func parseSeconds(s string) int {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || f <= 0 {
		return 0
	}
	return int(math.Round(f))
}

// annotateVideos заполняет VideoMeta у скачанных видео. Файлы ещё должны
// лежать в dir. Если ffprobe не справился, видео уходит без размеров, как раньше.
func (c *YtDlpClient) annotateVideos(ctx context.Context, dir string, resp *models.DownloadResponse) {
	if c.probeBinary == "" || resp == nil {
		return
	}
	probe := func(name string) models.VideoMeta {
		meta, err := probeVideo(ctx, c.probeBinary, filepath.Join(dir, name))
		if err != nil {
			log.Printf("probe video %s: %v", name, err)
		}
		return meta
	}
	if len(resp.Data) > 0 && resp.FileName != "" {
		resp.Video = probe(resp.FileName)
	}
	for i := range resp.Items {
		if resp.Items[i].Type == "video" {
			resp.Items[i].Video = probe(resp.Items[i].FileName)
		}
	}
}
