package logs

import (
	"os"
	"strings"
	"testing"
)

func TestAnalyze_EmptyFile(t *testing.T) {
	f, err := os.CreateTemp("", "test*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	res, err := Analyze(f.Name(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if res.TotalLines != 0 {
		t.Errorf("expected 0 lines, got %d", res.TotalLines)
	}
}

func TestAnalyze_WithErrors(t *testing.T) {
	f, err := os.CreateTemp("", "test*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("ERROR connection refused\nERROR timeout\nWARN retrying\n")
	f.Close()
	res, err := Analyze(f.Name(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if res.ErrorCount != 2 {
		t.Errorf("expected 2 errors, got %d", res.ErrorCount)
	}
	if res.WarningCount != 1 {
		t.Errorf("expected 1 warning, got %d", res.WarningCount)
	}
	md := res.ToMarkdown()
	if !strings.Contains(md, "connection refused") {
		t.Errorf("expected error cluster in markdown")
	}
}

func TestAnalyze_Focus(t *testing.T) {
	f, err := os.CreateTemp("", "test*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("ERROR db timeout\nINFO ping ok\nERROR db timeout\n")
	f.Close()
	res, err := Analyze(f.Name(), Options{Focus: "db"})
	if err != nil {
		t.Fatal(err)
	}
	if res.TotalLines != 2 {
		t.Errorf("expected 2 focused lines, got %d", res.TotalLines)
	}
}
