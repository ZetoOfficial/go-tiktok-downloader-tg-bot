package clients

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
)

func TestParseFFprobe(t *testing.T) {
	tests := []struct {
		name string
		json string
		want models.VideoMeta
	}{
		{
			name: "vertical",
			json: `{"streams":[{"width":720,"height":1280,"duration":"14.6"}],"format":{"duration":"14.7"}}`,
			want: models.VideoMeta{Width: 720, Height: 1280, Duration: 15},
		},
		{
			name: "rotated by side data",
			json: `{"streams":[{"width":1920,"height":1080,"side_data_list":[{"rotation":-90}]}],"format":{"duration":"3.2"}}`,
			want: models.VideoMeta{Width: 1080, Height: 1920, Duration: 3},
		},
		{
			name: "rotated by tag",
			json: `{"streams":[{"width":1280,"height":720,"tags":{"rotate":"270"},"duration":"5"}]}`,
			want: models.VideoMeta{Width: 720, Height: 1280, Duration: 5},
		},
		{
			name: "upside down keeps size",
			json: `{"streams":[{"width":1280,"height":720,"side_data_list":[{"rotation":180}],"duration":"5"}]}`,
			want: models.VideoMeta{Width: 1280, Height: 720, Duration: 5},
		},
		{
			name: "anamorphic pixels",
			json: `{"streams":[{"width":540,"height":960,"sample_aspect_ratio":"4:3","duration":"1"}]}`,
			want: models.VideoMeta{Width: 720, Height: 960, Duration: 1},
		},
		{
			name: "unknown sample aspect ratio",
			json: `{"streams":[{"width":720,"height":720,"sample_aspect_ratio":"0:1","duration":"N/A"}],"format":{"duration":"2.4"}}`,
			want: models.VideoMeta{Width: 720, Height: 720, Duration: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFFprobe([]byte(tt.json))
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseFFprobeErrors(t *testing.T) {
	for _, data := range []string{`not json`, `{"streams":[]}`, `{"streams":[{"width":0,"height":0}]}`} {
		if _, err := parseFFprobe([]byte(data)); err == nil {
			t.Errorf("parseFFprobe(%s) error = nil", data)
		}
	}
}

func TestProbeVideoRealFile(t *testing.T) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		t.Skip("ffmpeg not installed")
	}
	if _, err := exec.LookPath("ffprobe"); err != nil {
		t.Skip("ffprobe not installed")
	}
	path := filepath.Join(t.TempDir(), "clip.mp4")
	gen := exec.Command("ffmpeg", "-v", "error", "-f", "lavfi", "-i", "color=c=black:s=360x640:d=2",
		"-c:v", "libx264", "-pix_fmt", "yuv420p", path)
	if out, err := gen.CombinedOutput(); err != nil {
		t.Skipf("cannot generate test video: %v: %s", err, out)
	}
	got, err := probeVideo(context.Background(), "ffprobe", path)
	if err != nil {
		t.Fatal(err)
	}
	if want := (models.VideoMeta{Width: 360, Height: 640, Duration: 2}); got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestAnnotateVideosSkipsBrokenFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "clip.mp4", []byte("not a video"))
	resp := &models.DownloadResponse{FileName: "clip.mp4", Data: []byte("not a video")}
	c := &YtDlpClient{probeBinary: "ffprobe-missing-binary"}
	c.annotateVideos(context.Background(), dir, resp)
	if resp.Video != (models.VideoMeta{}) {
		t.Errorf("video meta = %+v, want zero", resp.Video)
	}
}
