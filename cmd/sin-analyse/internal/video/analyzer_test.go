package video

import (
	"strings"
	"testing"
)

func TestAnalyze_MissingFile(t *testing.T) {
	_, err := Analyze("/nonexistent.mp4", Options{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestToMarkdown(t *testing.T) {
	r := &Result{Path: "/v/test.mp4", SizeKB: 100}
	md := r.ToMarkdown()
	if !strings.Contains(md, "test.mp4") {
		t.Errorf("expected filename in markdown, got: %s", md)
	}
}
