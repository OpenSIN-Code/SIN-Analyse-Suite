package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Options struct {
	DDL        bool
	SampleSize int
}

type ColumnInfo struct {
	Name     string `json:"name"`
	TypeHint string `json:"type_hint"`
	NonNull  int    `json:"nonnull"`
	Nulls    int    `json:"nulls"`
}

type Result struct {
	Path        string            `json:"path"`
	SizeKB      int64             `json:"size_kb"`
	RowCount    int               `json:"row_count"`
	ColumnCount int               `json:"column_count"`
	Columns     []ColumnInfo      `json:"columns,omitempty"`
	SampleRows  []map[string]string `json:"sample_rows,omitempty"`
	DDL         string            `json:"ddl,omitempty"`
}

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
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read csv %s: %w", path, err)
	}
	if len(records) == 0 {
		return res, nil
	}
	headers := records[0]
	res.ColumnCount = len(headers)
	res.Columns = make([]ColumnInfo, len(headers))
	for i, h := range headers {
		res.Columns[i] = ColumnInfo{Name: h, TypeHint: "text"}
	}
	if len(records) > 1 {
		res.RowCount = len(records) - 1
		for _, row := range records[1:] {
			for i, val := range row {
				if i < len(res.Columns) {
					if val == "" {
						res.Columns[i].Nulls++
					} else {
						res.Columns[i].NonNull++
					}
				}
			}
		}
		sampleSize := opts.SampleSize
		if sampleSize <= 0 {
			sampleSize = 5
		}
		for i := 1; i < len(records) && i <= sampleSize; i++ {
			row := make(map[string]string)
			for j, h := range headers {
				if j < len(records[i]) {
					row[h] = records[i][j]
				}
			}
			res.SampleRows = append(res.SampleRows, row)
		}
		if opts.DDL {
			res.DDL = generateDDL(path, res.Columns)
		}
	}
	return res, nil
}

func generateDDL(path string, cols []ColumnInfo) string {
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, name)
	if name == "" {
		name = "imported_data"
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", name))
	for i, c := range cols {
		sqlType := "TEXT"
		hint := strings.ToLower(c.TypeHint)
		if strings.Contains(hint, "int") || strings.Contains(hint, "num") {
			sqlType = "INTEGER"
		} else if strings.Contains(hint, "float") || strings.Contains(hint, "real") {
			sqlType = "REAL"
		} else if strings.Contains(hint, "bool") {
			sqlType = "BOOLEAN"
		}
		nullable := "NOT NULL"
		if c.Nulls > 0 {
			nullable = ""
		}
		comma := ","
		if i == len(cols)-1 {
			comma = ""
		}
		b.WriteString(fmt.Sprintf("  %s %s %s%s\n", sanitizeName(c.Name), sqlType, nullable, comma))
	}
	b.WriteString(");")
	return b.String()
}

func sanitizeName(name string) string {
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

func (r *Result) ToMarkdown() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("## Data Analysis: %s\n\n", filepath.Base(r.Path)))
	b.WriteString(fmt.Sprintf("- **Rows**: %d\n", r.RowCount))
	b.WriteString(fmt.Sprintf("- **Columns**: %d\n", r.ColumnCount))
	if len(r.Columns) > 0 {
		b.WriteString("\n### Columns\n\n| Name | Type | Non-Null | Nulls |\n|------|------|----------|-------|\n")
		for _, c := range r.Columns {
			b.WriteString(fmt.Sprintf("| %s | %s | %d | %d |\n", c.Name, c.TypeHint, c.NonNull, c.Nulls))
		}
	}
	if len(r.SampleRows) > 0 {
		b.WriteString("\n### Sample Rows\n\n```\n")
		for _, row := range r.SampleRows {
			vals := make([]string, 0, len(row))
			for _, v := range row {
				vals = append(vals, v)
			}
			b.WriteString(strings.Join(vals, ", ") + "\n")
		}
		b.WriteString("```\n")
	}
	if r.DDL != "" {
		b.WriteString("\n### DDL\n\n```sql\n" + r.DDL + "\n```\n")
	}
	return b.String()
}
