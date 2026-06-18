package image

import (
	"os"
	"strings"
	"testing"
)

func TestAnalyze_NonExistent(t *testing.T) {
	_, err := Analyze("/nonexistent.png", Options{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAnalyze_ToMarkdown(t *testing.T) {
	f, err := os.CreateTemp("", "test*.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.Write([]byte("fakejpegdata"))
	f.Close()
	res, err := Analyze(f.Name(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	md := res.ToMarkdown()
	if !strings.Contains(md, "test") {
		t.Errorf("expected filename in markdown, got: %s", md)
	}
}
