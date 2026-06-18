# SIN-Analyse-Suite

Multimodal preprocessing pipelines for SIN-Code — image, video, PDF, logs, data, and audio analysis via MCP protocol.

## Features

| Tool | Description | External Dep |
|---|---|---|
| `analyse_image` | Metadata, resize, OCR | Tesseract |
| `analyse_video` | Keyframes, transcription | ffmpeg + Whisper |
| `analyse_pdf` | Text, tables, OCR | pdfcpu + Tesseract |
| `analyse_logs` | Error clustering, filtering | — |
| `analyse_data` | Schema, stats, DDL | — |
| `analyse_audio` | Transcription, diarization | Whisper |

## Quick Start

```bash
go install github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse@latest

sin-analyse image screenshot.png
sin-analyse video demo.mp4
sin-analyse logs /var/log/syslog --focus error
sin-analyse data users.csv --ddl
sin-analyse audio meeting.mp3
sin-analyse serve   # MCP server
```

### Dependencies

```bash
brew install ffmpeg tesseract        # macOS
sudo apt install ffmpeg tesseract-ocr # Ubuntu
```

## Integration with SIN-Code

```bash
# Register as MCP server
sin-code mcp register sin-analyse "$(which sin-analyse)" serve
```

## Development

```bash
make build    # Build binary
make test     # Run tests with race detection
make check    # Build + lint + test
```

## License

MIT
