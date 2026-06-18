package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRegisterOpencode(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "opencode.json")

	err := registerOpencode(configPath, "/fake/sin-analyse")
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}

	mcp, ok := config["mcp"].(map[string]any)
	if !ok {
		t.Fatal("mcp key missing")
	}

	entry, ok := mcp["sin-analyse"].(map[string]any)
	if !ok {
		t.Fatal("sin-analyse entry missing")
	}

	cmd, ok := entry["command"].([]any)
	if !ok || len(cmd) != 2 || cmd[0] != "/fake/sin-analyse" || cmd[1] != "serve" {
		t.Fatalf("unexpected command: %v", entry["command"])
	}

	if entry["enabled"] != true {
		t.Error("expected enabled=true")
	}

	if entry["type"] != "local" {
		t.Error("expected type=local")
	}
}

func TestRegisterClaude(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "mcp.json")

	err := registerClaude(configPath, "/fake/sin-analyse")
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}

	servers, ok := config["mcpServers"].(map[string]any)
	if !ok {
		t.Fatal("mcpServers key missing")
	}

	entry, ok := servers["sin-analyse"].(map[string]any)
	if !ok {
		t.Fatal("sin-analyse entry missing")
	}

	if entry["command"] != "/fake/sin-analyse" {
		t.Fatalf("unexpected command: %v", entry["command"])
	}

	args, ok := entry["args"].([]any)
	if !ok || len(args) != 1 || args[0] != "serve" {
		t.Fatalf("unexpected args: %v", entry["args"])
	}
}

func TestRegisterIdempotent(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "opencode.json")

	if err := registerOpencode(configPath, "/fake/sin-analyse"); err != nil {
		t.Fatal(err)
	}
	if err := registerOpencode(configPath, "/updated/sin-analyse"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}

	mcp := config["mcp"].(map[string]any)
	entry := mcp["sin-analyse"].(map[string]any)
	cmd := entry["command"].([]any)

	if cmd[0] != "/updated/sin-analyse" {
		t.Fatalf("expected updated path, got %v", cmd[0])
	}
}

func TestRegisterPreservesExisting(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "opencode.json")

	existing := map[string]any{
		"mcp": map[string]any{
			"other-server": map[string]any{
				"command": []string{"other-bin"},
				"enabled": true,
				"type":    "local",
			},
		},
	}
	data, _ := json.MarshalIndent(existing, "", "  ")
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := registerOpencode(configPath, "/fake/sin-analyse"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}

	mcp := config["mcp"].(map[string]any)
	if _, ok := mcp["other-server"]; !ok {
		t.Error("existing server was removed")
	}
	if _, ok := mcp["sin-analyse"]; !ok {
		t.Error("sin-analyse was not added")
	}
}

func TestFindTarget(t *testing.T) {
	target, err := findTarget("opencode")
	if err != nil {
		t.Fatal(err)
	}
	if target.Name != "opencode" {
		t.Fatalf("expected opencode, got %s", target.Name)
	}

	_, err = findTarget("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}
