//go:build e2e

package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "sin-analyse")
	cmd := exec.Command("go", "build", "-o", bin, "./cmd/sin-analyse")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %s\n%s", err, out)
	}
	return bin
}

func TestCLI_Image(t *testing.T) {
	bin := buildBinary(t)
	tmp := t.TempDir()
	img := filepath.Join(tmp, "test.jpg")
	os.WriteFile(img, []byte("fakejpeg"), 0644)

	out, err := exec.Command(bin, "image", img, "--json").CombinedOutput()
	if err != nil {
		t.Fatalf("cli: %s\n%s", err, out)
	}
	var res map[string]interface{}
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatalf("json: %s\n%s", err, out)
	}
	if res["path"] != img {
		t.Errorf("expected path %s, got %v", img, res["path"])
	}
}

func TestCLI_Logs(t *testing.T) {
	bin := buildBinary(t)
	tmp := t.TempDir()
	log := filepath.Join(tmp, "test.log")
	os.WriteFile(log, []byte("ERROR test error\nINFO ok\n"), 0644)

	out, err := exec.Command(bin, "logs", log, "--json").CombinedOutput()
	if err != nil {
		t.Fatalf("cli: %s\n%s", err, out)
	}
	var res map[string]interface{}
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatalf("json: %s\n%s", err, out)
	}
	if res["total_lines"] != 2.0 {
		t.Errorf("expected 2 lines, got %v", res["total_lines"])
	}
}

func TestCLI_Data(t *testing.T) {
	bin := buildBinary(t)
	tmp := t.TempDir()
	csv := filepath.Join(tmp, "test.csv")
	os.WriteFile(csv, []byte("name,age\nAlice,30\nBob,25\n"), 0644)

	out, err := exec.Command(bin, "data", csv, "--ddl", "--json").CombinedOutput()
	if err != nil {
		t.Fatalf("cli: %s\n%s", err, out)
	}
	var res map[string]interface{}
	if err := json.Unmarshal(out, &res); err != nil {
		t.Fatalf("json: %s\n%s", err, out)
	}
	if res["ddl"] == nil {
		t.Error("expected DDL in JSON output")
	}
}
