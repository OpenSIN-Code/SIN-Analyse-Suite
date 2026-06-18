package audio

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Options struct {
	Language string
}

type Result struct {
	Path          string  `json:"path"`
	SizeKB        int64   `json:"size_kb"`
	DurationSec   float64 `json:"duration_sec"`
	Format        string  `json:"format"`
	SampleRate    int     `json:"sample_rate"`
	Channels      int     `json:"channels"`
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
	res := &Result{Path: path, SizeKB: info.Size() / 1024, Format: strings.TrimPrefix(filepath.Ext(path), ".")}
	dur, _ := ffprobeDuration(path)
	res.DurationSec = dur
	sr, ch, _ := ffprobeAudioInfo(path)
	res.SampleRate = sr
	res.Channels = ch
	trans, _ := transcribe(path, opts.Language)
	res.Transcription = trans
	return res, nil
}

func ffprobeDuration(path string) (float64, error) {
	out, err := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "csv=p=0", path).Output()
	if err != nil {
		return 0, err
	}
	var d float64
	if _, err := fmt.Sscanf(strings.TrimSpace(string(out)), "%f", &d); err != nil {
		return 0, fmt.Errorf("parse duration: %w", err)
	}
	return d, nil
}

func ffprobeAudioInfo(path string) (sr, ch int, err error) {
	out, err := exec.Command("ffprobe", "-v", "error", "-select_streams", "a:0", "-show_entries", "stream=sample_rate,channels", "-of", "csv=p=0", path).Output()
	if err != nil {
		return 0, 0, err
	}
	parts := strings.Split(strings.TrimSpace(string(out)), ",")
	if len(parts) >= 2 {
		if _, err := fmt.Sscanf(parts[0], "%d", &sr); err != nil {
			return 0, 0, fmt.Errorf("parse sample_rate: %w", err)
		}
		if _, err := fmt.Sscanf(parts[1], "%d", &ch); err != nil {
			return 0, 0, fmt.Errorf("parse channels: %w", err)
		}
	}
	return
}

func (r *Result) ToMarkdown() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("## Audio Analysis: %s\n\n", filepath.Base(r.Path)))
	b.WriteString(fmt.Sprintf("- **Size**: %d KB\n", r.SizeKB))
	if r.DurationSec > 0 {
		b.WriteString(fmt.Sprintf("- **Duration**: %.1f sec\n", r.DurationSec))
	}
	b.WriteString(fmt.Sprintf("- **Format**: %s\n", r.Format))
	if r.SampleRate > 0 {
		b.WriteString(fmt.Sprintf("- **Sample Rate**: %d Hz\n", r.SampleRate))
	}
	if r.Transcription != "" {
		preview := r.Transcription
		if len(preview) > 300 {
			preview = preview[:300] + "..."
		}
		b.WriteString(fmt.Sprintf("- **Transcript**: %s\n", preview))
	}
	return b.String()
}
