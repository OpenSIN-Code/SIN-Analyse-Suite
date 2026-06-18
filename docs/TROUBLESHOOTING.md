# Troubleshooting

## ffmpeg/tesseract not found
```bash
brew install ffmpeg tesseract        # macOS
sudo apt install ffmpeg tesseract-ocr # Ubuntu
```

## No whisper backend
Set `OPENAI_API_KEY` or install whisper.cpp.

## Permission blocked (SIN-Code)
```bash
sin-code permission add analyse__* allow
```
