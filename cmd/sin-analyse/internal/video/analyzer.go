package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Options struct {
	KeyframeInterval int
}

type Result struct {
	Path          string  `json:"path"`
	SizeKB        int64   `json:"size_kb"`
	DurationSec   float64 `json:"duration_sec"`
	Codec         string  `json:"codec"`
	KeyframeCount int     `json:"keyframes,omitempty"`
	Transcription string  `json:"transcription,omitempty"`
}

func Analyze(path string, opts Options) (*Result, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return nil, fmt.Errorf("ffmpeg required but not found on PATH")
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}
	res := &Result{Path: path, SizeKB: info.Size() / 1024}
	dur, _ := getDuration(path)
	res.DurationSec = dur
	codec, _ := getCodec(path)
	res.Codec = codec
	return res, nil
}

func getDuration(path string) (float64, error) {
	out, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "csv=p=0", path).Output()
	if err != nil {
		return 0, err
	}
	var d float64
	fmt.Sscanf(strings.TrimSpace(string(out)), "%f", &d)
	return d, nil
}

func getCodec(path string) (string, error) {
	out, err := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=codec_name", "-of", "csv=p=0", path).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (r *Result) ToMarkdown() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("## Video Analysis: %s\n\n", filepath.Base(r.Path)))
	b.WriteString(fmt.Sprintf("- **Path**: `%s`\n", r.Path))
	b.WriteString(fmt.Sprintf("- **Size**: %d KB\n", r.SizeKB))
	if r.DurationSec > 0 {
		b.WriteString(fmt.Sprintf("- **Duration**: %.1f sec\n", r.DurationSec))
	}
	if r.Codec != "" {
		b.WriteString(fmt.Sprintf("- **Codec**: %s\n", r.Codec))
	}
	return b.String()
}
