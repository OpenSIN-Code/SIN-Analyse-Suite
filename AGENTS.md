# AGENTS.md — SIN-Analyse-Suite

## Purpose
Go-native MCP skill server providing multimodal preprocessing pipelines for SIN-Code.

## Architecture
- **M2**: Single static Go binary, CGO_ENABLED=0
- **M6**: SIN tools over naive built-ins
- **M7**: `go test -race` mandatory
- Read-only — never modifies input files
- ffmpeg/whisper/Tesseract called as subprocesses (Bridged-External)

## Change Policy
- Run `make check` before commit
- Conventional commits: `feat(pkg):`, `fix(pkg):`, `docs:`, `test:`
- New analyzer: add `internal/` package, register in `mcp/server.go`, add CLI command

## Testing Gates
```bash
go build ./...
go test ./... -race -count=1
