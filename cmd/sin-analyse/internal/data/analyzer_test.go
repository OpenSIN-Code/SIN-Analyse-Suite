package data

import (
	"os"
	"strings"
	"testing"
)

func TestAnalyze_CSV(t *testing.T) {
	f, err := os.CreateTemp("", "test*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString("name,age,email\nAlice,30,alice@test.com\nBob,25,bob@test.com\n"); err != nil {
		t.Fatal(err)
	}
	f.Close()
	res, err := Analyze(f.Name(), Options{SampleSize: 2, DDL: true})
	if err != nil {
		t.Fatal(err)
	}
	if res.RowCount != 2 {
		t.Errorf("expected 2 rows, got %d", res.RowCount)
	}
	if res.ColumnCount != 3 {
		t.Errorf("expected 3 columns, got %d", res.ColumnCount)
	}
	if len(res.SampleRows) != 2 {
		t.Errorf("expected 2 sample rows, got %d", len(res.SampleRows))
	}
	if res.DDL == "" {
		t.Error("expected DDL")
	}
	md := res.ToMarkdown()
	if !strings.Contains(md, "CREATE TABLE") {
		t.Errorf("expected DDL in markdown, got: %s", md)
	}
}

func TestAnalyze_Empty(t *testing.T) {
	f, err := os.CreateTemp("", "empty*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	res, err := Analyze(f.Name(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	if res.RowCount != 0 {
		t.Errorf("expected 0 rows, got %d", res.RowCount)
	}
}

func TestAnalyze_NonExistent(t *testing.T) {
	_, err := Analyze("/nonexistent.csv", Options{})
	if err == nil {
		t.Fatal("expected error")
	}
}
