package audio

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// sin-debt: perf, upgrade: replace with pure-Go whisper bindings (e.g. whspr)

func transcribe(path, language string) (string, error) {
	if _, err := exec.LookPath("whisper"); err == nil {
		return transcribeWhisperCpp(path, language)
	}
	if _, err := exec.LookPath("python3"); err == nil {
		if text, err := transcribePythonWhisper(path, language); err == nil {
			return text, nil
		}
	}
	if os.Getenv("OPENAI_API_KEY") != "" {
		return transcribeOpenAI(path, language)
	}
	return "", fmt.Errorf("no whisper backend — install whisper.cpp, python -m pip install openai-whisper, or set OPENAI_API_KEY")
}

// sin-debt: portability, upgrade: ship static whisper.cpp binary in release tarball
func transcribeWhisperCpp(path, language string) (string, error) {
	args := []string{"--model", "base", "--output-txt", path}
	if language != "" && language != "auto" {
		args = append(args, "--language", language)
	}
	out, err := exec.Command("whisper", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("whisper.cpp: %s: %w", string(out), err)
	}
	data, err := os.ReadFile(path + ".txt")
	if err != nil {
		return "", fmt.Errorf("read whisper output: %w", err)
	}
	return string(data), nil
}

// sin-debt: latency, upgrade: switch to faster-whisper or distil-whisper
func transcribePythonWhisper(path, language string) (string, error) {
	script := `import sys, whisper
model = whisper.load_model("base")
lang = None if len(sys.argv) < 3 else sys.argv[2]
result = model.transcribe(sys.argv[1], language=lang)
print(result["text"])
`
	args := []string{"-c", script, "--", path}
	if language != "" && language != "auto" {
		args = append(args, language)
	}
	cmd := exec.Command("python3", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("python-whisper: %s: %w", stderr.String(), err)
	}
	return string(out), nil
}

// sin-debt: api-cost, upgrade: use batch endpoint for multi-file transcription
func transcribeOpenAI(path, language string) (string, error) {
	return "", fmt.Errorf("OpenAI Whisper API not wired — implement HTTP multipart POST to https://api.openai.com/v1/audio/transcriptions")
}
