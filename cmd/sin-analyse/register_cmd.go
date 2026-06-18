package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type AgentTarget struct {
	Name        string
	DisplayName string
	ConfigPath  func(home string) string
	Format      string
}

var agentTargets = []AgentTarget{
	{
		Name:        "opencode",
		DisplayName: "opencode CLI",
		ConfigPath:  func(home string) string { return filepath.Join(home, ".config", "opencode", "opencode.json") },
		Format:      "opencode",
	},
	{
		Name:        "claude-code",
		DisplayName: "Claude Code",
		ConfigPath:  func(home string) string { return filepath.Join(home, ".claude", "mcp.json") },
		Format:      "claude",
	},
	{
		Name:        "codex",
		DisplayName: "Codex CLI",
		ConfigPath:  func(home string) string { return filepath.Join(home, ".codex", "config.json") },
		Format:      "codex",
	},
	{
		Name:        "cursor",
		DisplayName: "Cursor",
		ConfigPath:  func(home string) string { return filepath.Join(home, ".cursor", "mcp.json") },
		Format:      "claude",
	},
	{
		Name:        "windsurf",
		DisplayName: "Windsurf",
		ConfigPath:  func(home string) string { return filepath.Join(home, ".codeium", "windsurf", "mcp_config.json") },
		Format:      "claude",
	},
	{
		Name:        "cline",
		DisplayName: "Cline",
		ConfigPath: func(home string) string {
			return filepath.Join(home, ".vscode", "extensions", "saoudrizwan.claude-dev", "settings", "cline_mcp_settings.json")
		},
		Format: "claude",
	},
}

var registerCmd = &cobra.Command{
	Use:   "register --agent <target> [--path /custom/binary]",
	Short: "Register sin-analyse MCP server in an agent's config",
	Long: `Register the sin-analyse MCP server in an external agent's configuration file.

Supported agents:
  opencode      ~/.config/opencode/opencode.json
  claude-code   ~/.claude/mcp.json
  codex         ~/.codex/config.json
  cursor        ~/.cursor/mcp.json
  windsurf      ~/.codeium/windsurf/mcp_config.json
  cline         ~/.vscode/extensions/saoudrizwan.claude-dev/settings/cline_mcp_settings.json

The --path flag overrides the binary path (default: auto-detect from PATH or current binary).`,
	RunE: runRegister,
}

var registerAgent string
var registerPath string

func init() {
	registerCmd.Flags().StringVar(&registerAgent, "agent", "", "Target agent (opencode|claude-code|codex|cursor|windsurf|cline)")
	registerCmd.Flags().StringVar(&registerPath, "path", "", "Override binary path (default: auto-detect)")
	rootCmd.AddCommand(registerCmd)
}

func runRegister(cmd *cobra.Command, args []string) error {
	if registerAgent == "" {
		return fmt.Errorf("--agent is required (opencode|claude-code|codex|cursor|windsurf|cline)")
	}

	target, err := findTarget(registerAgent)
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home dir: %w", err)
	}

	binaryPath := registerPath
	if binaryPath == "" {
		binaryPath = autoDetectBinary()
	}
	if binaryPath == "" {
		return fmt.Errorf("could not find sin-analyse binary — use --path to specify")
	}

	configPath := target.ConfigPath(home)

	switch target.Format {
	case "opencode":
		return registerOpencode(configPath, binaryPath)
	case "claude":
		return registerClaude(configPath, binaryPath)
	case "codex":
		return registerCodex(configPath, binaryPath)
	default:
		return fmt.Errorf("unknown format: %s", target.Format)
	}
}

func findTarget(name string) (*AgentTarget, error) {
	for i := range agentTargets {
		if agentTargets[i].Name == name {
			return &agentTargets[i], nil
		}
	}
	return nil, fmt.Errorf("unknown agent: %s (valid: opencode|claude-code|codex|cursor|windsurf|cline)", name)
}

func autoDetectBinary() string {
	if exe, err := os.Executable(); err == nil {
		return exe
	}
	for _, p := range strings.Split(os.Getenv("PATH"), ":") {
		candidate := filepath.Join(p, "sin-analyse")
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}
	}
	return ""
}

func registerOpencode(configPath, binaryPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return writeJSON(configPath, map[string]any{
				"mcp": map[string]any{
					"sin-analyse": map[string]any{
						"command":     []string{binaryPath, "serve"},
						"enabled":     true,
						"type":        "local",
						"description": "SIN-Analyse-Suite — multimodal preprocessing (image, video, PDF, logs, data, audio)",
					},
				},
			})
		}
		return fmt.Errorf("read %s: %w", configPath, err)
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parse %s: %w", configPath, err)
	}

	mcp, ok := config["mcp"].(map[string]any)
	if !ok {
		mcp = map[string]any{}
		config["mcp"] = mcp
	}

	mcp["sin-analyse"] = map[string]any{
		"command":     []string{binaryPath, "serve"},
		"enabled":     true,
		"type":        "local",
		"description": "SIN-Analyse-Suite — multimodal preprocessing (image, video, PDF, logs, data, audio)",
	}

	return writeJSON(configPath, config)
}

func registerClaude(configPath, binaryPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return writeJSON(configPath, map[string]any{
				"mcpServers": map[string]any{
					"sin-analyse": map[string]any{
						"command": binaryPath,
						"args":    []string{"serve"},
					},
				},
			})
		}
		return fmt.Errorf("read %s: %w", configPath, err)
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parse %s: %w", configPath, err)
	}

	servers, ok := config["mcpServers"].(map[string]any)
	if !ok {
		servers = map[string]any{}
		config["mcpServers"] = servers
	}

	servers["sin-analyse"] = map[string]any{
		"command": binaryPath,
		"args":    []string{"serve"},
	}

	return writeJSON(configPath, config)
}

func registerCodex(configPath, binaryPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return writeJSON(configPath, map[string]any{
				"mcpServers": map[string]any{
					"sin-analyse": map[string]any{
						"command":     binaryPath,
						"args":        []string{"serve"},
						"description": "SIN-Analyse-Suite — multimodal preprocessing",
					},
				},
			})
		}
		return fmt.Errorf("read %s: %w", configPath, err)
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parse %s: %w", configPath, err)
	}

	servers, ok := config["mcpServers"].(map[string]any)
	if !ok {
		servers = map[string]any{}
		config["mcpServers"] = servers
	}

	servers["sin-analyse"] = map[string]any{
		"command":     binaryPath,
		"args":        []string{"serve"},
		"description": "SIN-Analyse-Suite — multimodal preprocessing",
	}

	return writeJSON(configPath, config)
}

func writeJSON(path string, v any) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	fmt.Printf("Registered sin-analyse in %s\n", path)
	return nil
}
