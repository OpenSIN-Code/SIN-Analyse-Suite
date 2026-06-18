package pdf

import (
	"os"
	"strings"
	"testing"
)

func TestAnalyze_NonExistent(t *testing.T) {
	_, err := Analyze("/nonexistent.pdf", Options{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestAnalyze_ToMarkdown(t *testing.T) {
	f, err := os.CreateTemp("", "test*.pdf")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.Write([]byte("%PDF-1.4 fake")); err != nil {
		t.Fatal(err)
	}
	f.Close()
	res, err := Analyze(f.Name(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	md := res.ToMarkdown()
	if !strings.Contains(md, "test") {
		t.Errorf("expected filename, got: %s", md)
	}
}
