package audio

import (
	"strings"
	"testing"
)

func TestAnalyze_MissingFile(t *testing.T) {
	_, err := Analyze("/nonexistent.mp3", Options{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestToMarkdown(t *testing.T) {
	r := &Result{Path: "/a/m.mp3", SizeKB: 1024, Format: "mp3", Transcription: "hello test"}
	md := r.ToMarkdown()
	if !strings.Contains(md, "m.mp3") {
		t.Errorf("expected filename, got: %s", md)
	}
	if !strings.Contains(md, "hello test") {
		t.Errorf("expected transcript, got: %s", md)
	}
}
