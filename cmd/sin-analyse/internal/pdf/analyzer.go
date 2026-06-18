package pdf

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
	Path       string `json:"path"`
	SizeKB     int64  `json:"size_kb"`
	PageCount  int    `json:"page_count"`
	Text       string `json:"text,omitempty"`
	OCREnabled bool   `json:"ocr_enabled"`
}

func Analyze(path string, opts Options) (*Result, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}
	res := &Result{Path: path, SizeKB: info.Size() / 1024, OCREnabled: opts.OCR}
	pages, text, err := extractText(path)
	if err == nil {
		res.PageCount = pages
		res.Text = text
	}
	return res, nil
}

// sin-debt: portability, upgrade: vendor pdfcpu as direct Go dependency instead of subprocess
func extractText(path string) (int, string, error) {
	return 0, "", fmt.Errorf("pdfcpu extraction not wired — add github.com/pdfcpu/pdfcpu dependency and wire pdfcpu.Extract")
}

func (r *Result) ToMarkdown() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("## PDF Analysis: %s\n\n", filepath.Base(r.Path)))
	b.WriteString(fmt.Sprintf("- **Path**: `%s`\n", r.Path))
	b.WriteString(fmt.Sprintf("- **Size**: %d KB\n", r.SizeKB))
	if r.PageCount > 0 {
		b.WriteString(fmt.Sprintf("- **Pages**: %d\n", r.PageCount))
	}
	if r.Text != "" {
		preview := r.Text
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}
		b.WriteString(fmt.Sprintf("- **Preview**: %s\n", preview))
	}
	return b.String()
}
