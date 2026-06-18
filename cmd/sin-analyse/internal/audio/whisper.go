package audio

import (
	"fmt"
	"os"
	"os/exec"
)

func transcribe(path, language string) (string, error) {
	if _, err := exec.LookPath("whisper"); err == nil {
		return transcribeWhisperCpp(path, language)
	}
	if os.Getenv("OPENAI_API_KEY") != "" {
		return transcribeOpenAI(path, language)
	}
	return "", fmt.Errorf("no whisper backend; install whisper.cpp or set OPENAI_API_KEY")
}

func transcribeWhisperCpp(path, language string) (string, error) {
	args := []string{"--model", "base", "--output-txt", path}
	if language != "" && language != "auto" {
		args = append(args, "--language", language)
	}
	out, err := exec.Command("whisper", args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("whisper: %s: %w", string(out), err)
	}
	data, err := os.ReadFile(path + ".txt")
	if err != nil {
		return "", fmt.Errorf("read whisper output: %w", err)
	}
	return string(data), nil
}

func transcribeOpenAI(path, language string) (string, error) {
	return "", fmt.Errorf("OpenAI Whisper API not wired")
}
