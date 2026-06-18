package logs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Options struct {
	Focus string
	Since string
}

type Result struct {
	Path          string         `json:"path"`
	SizeKB        int64          `json:"size_kb"`
	TotalLines    int            `json:"total_lines"`
	ErrorCount    int            `json:"error_count"`
	WarningCount  int            `json:"warning_count"`
	ErrorClusters []ErrorCluster `json:"error_clusters,omitempty"`
	ErrorPatterns map[string]int `json:"error_patterns,omitempty"`
}

type ErrorCluster struct {
	Pattern string `json:"pattern"`
	Count   int    `json:"count"`
	Example string `json:"example"`
}

var errorRe = regexp.MustCompile(`(?i)\berror\b`)
var warningRe = regexp.MustCompile(`(?i)\bwarn(ing)?\b`)

func Analyze(path string, opts Options) (*Result, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}
	res := &Result{Path: path, SizeKB: info.Size() / 1024}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	patterns := make(map[string]int)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if opts.Focus != "" && !strings.Contains(line, opts.Focus) {
			continue
		}
		res.TotalLines++
		if errorRe.MatchString(line) {
			res.ErrorCount++
			patterns[normalizeErrorLine(line)]++
		}
		if warningRe.MatchString(line) {
			res.WarningCount++
		}
	}
	res.ErrorPatterns = patterns
	for pattern, cnt := range patterns {
		res.ErrorClusters = append(res.ErrorClusters, ErrorCluster{Pattern: pattern, Count: cnt})
	}
	return res, scanner.Err()
}

func normalizeErrorLine(line string) string {
	re := regexp.MustCompile(`\s+\d{4}[-/]\d{2}[-/]\d{2}|T\d{2}:\d{2}:\d{2}|\.go:\d+`)
	line = re.ReplaceAllString(line, " <ts>")
	if len(line) > 120 {
		line = line[:120]
	}
	return line
}

func (r *Result) ToMarkdown() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("## Log Analysis: %s\n\n", filepath.Base(r.Path)))
	b.WriteString(fmt.Sprintf("- **Lines**: %d\n", r.TotalLines))
	b.WriteString(fmt.Sprintf("- **Errors**: %d\n", r.ErrorCount))
	b.WriteString(fmt.Sprintf("- **Warnings**: %d\n", r.WarningCount))
	if len(r.ErrorClusters) > 0 {
		b.WriteString("\n### Error Clusters\n\n")
		for _, c := range r.ErrorClusters {
			b.WriteString(fmt.Sprintf("- `%s` — %dx\n", c.Pattern, c.Count))
		}
	}
	return b.String()
}
