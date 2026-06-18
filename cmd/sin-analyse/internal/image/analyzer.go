package image

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	OCR bool
}

type Result struct {
	Path    string `json:"path"`
	SizeKB  int64  `json:"size_kb"`
	Format  string `json:"format"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	OCRText string `json:"ocr_text,omitempty"`
}

func Analyze(path string, opts Options) (*Result, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}
	ext := strings.ToLower(filepath.Ext(path))
	res := &Result{Path: path, SizeKB: info.Size() / 1024, Format: ext}
	if opts.OCR {
		text, err := runOCR(path)
		if err != nil {
			return nil, fmt.Errorf("ocr: %w", err)
		}
		res.OCRText = text
	}
	return res, nil
}

func runOCR(path string) (string, error) {
	return "", fmt.Errorf("Tesseract not available")
}

func (r *Result) ToMarkdown() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("## Image Analysis: %s\n\n", filepath.Base(r.Path)))
	b.WriteString(fmt.Sprintf("- **Path**: `%s`\n", r.Path))
	b.WriteString(fmt.Sprintf("- **Size**: %d KB\n", r.SizeKB))
	b.WriteString(fmt.Sprintf("- **Format**: %s\n", r.Format))
	if r.OCRText != "" {
		b.WriteString(fmt.Sprintf("- **OCR Text**: %s\n", r.OCRText))
	}
	return b.String()
}
