package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/audio"
	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/data"
	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/image"
	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/logs"
	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/pdf"
	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/video"
)

func RegisterAll(s *server.MCPServer) error {
	tools := []struct {
		tool    mcp.Tool
		handler server.ToolHandlerFunc
	}{
		{
			tool: mcp.NewTool("analyse_image",
				mcp.WithDescription("Analyse an image file: metadata, OCR"),
				mcp.WithString("path", mcp.Required(), mcp.Description("Path to image file")),
				mcp.WithBoolean("ocr", mcp.Description("Run OCR via Tesseract")),
			),
			handler: handleAnalyseImage,
		},
		{
			tool: mcp.NewTool("analyse_video",
				mcp.WithDescription("Analyse a video file: keyframes, transcription"),
				mcp.WithString("path", mcp.Required(), mcp.Description("Path to video file")),
				mcp.WithNumber("keyframe_interval", mcp.Description("Keyframe interval (frames)")),
			),
			handler: handleAnalyseVideo,
		},
		{
			tool: mcp.NewTool("analyse_pdf",
				mcp.WithDescription("Analyse a PDF file: text, tables, OCR"),
				mcp.WithString("path", mcp.Required(), mcp.Description("Path to PDF file")),
				mcp.WithBoolean("ocr", mcp.Description("Run OCR for scanned pages")),
			),
			handler: handleAnalysePDF,
		},
		{
			tool: mcp.NewTool("analyse_logs",
				mcp.WithDescription("Analyse a log file: error clustering, filtering"),
				mcp.WithString("path", mcp.Required(), mcp.Description("Path to log file")),
				mcp.WithString("focus", mcp.Description("Focus on specific keyword")),
				mcp.WithString("since", mcp.Description("Filter entries since timestamp")),
			),
			handler: handleAnalyseLogs,
		},
		{
			tool: mcp.NewTool("analyse_data",
				mcp.WithDescription("Analyse a data file (CSV): schema, stats, DDL"),
				mcp.WithString("path", mcp.Required(), mcp.Description("Path to CSV file")),
				mcp.WithBoolean("ddl", mcp.Description("Generate CREATE TABLE DDL")),
				mcp.WithNumber("sample_size", mcp.Description("Number of sample rows")),
			),
			handler: handleAnalyseData,
		},
		{
			tool: mcp.NewTool("analyse_audio",
				mcp.WithDescription("Analyse an audio file: transcription, diarization"),
				mcp.WithString("path", mcp.Required(), mcp.Description("Path to audio file")),
				mcp.WithString("language", mcp.Description("Language (auto, en, de, fr)")),
			),
			handler: handleAnalyseAudio,
		},
	}
	for _, t := range tools {
		s.AddTool(t.tool, t.handler)
	}
	return nil
}

func resultJSON(v interface{}) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("marshal: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func handleAnalyseImage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseString(req, "path", "")
	if path == "" {
		return mcp.NewToolResultError("'path' is required"), nil
	}
	res, err := image.Analyze(path, image.Options{OCR: 	mcp.ParseBoolean(req, "ocr", false)})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return resultJSON(res)
}

func handleAnalyseVideo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseString(req, "path", "")
	if path == "" {
		return mcp.NewToolResultError("'path' is required"), nil
	}
	res, err := video.Analyze(path, video.Options{KeyframeInterval: mcp.ParseInt(req, "keyframe_interval", 30)})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return resultJSON(res)
}

func handleAnalysePDF(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseString(req, "path", "")
	if path == "" {
		return mcp.NewToolResultError("'path' is required"), nil
	}
	res, err := pdf.Analyze(path, pdf.Options{OCR: 	mcp.ParseBoolean(req, "ocr", false)})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return resultJSON(res)
}

func handleAnalyseLogs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseString(req, "path", "")
	if path == "" {
		return mcp.NewToolResultError("'path' is required"), nil
	}
	res, err := logs.Analyze(path, logs.Options{
		Focus: mcp.ParseString(req, "focus", ""),
		Since: mcp.ParseString(req, "since", ""),
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return resultJSON(res)
}

func handleAnalyseData(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseString(req, "path", "")
	if path == "" {
		return mcp.NewToolResultError("'path' is required"), nil
	}
	res, err := data.Analyze(path, data.Options{
		DDL:        mcp.ParseBoolean(req, "ddl", false),
		SampleSize: mcp.ParseInt(req, "sample_size", 5),
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return resultJSON(res)
}

func handleAnalyseAudio(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path := mcp.ParseString(req, "path", "")
	if path == "" {
		return mcp.NewToolResultError("'path' is required"), nil
	}
	res, err := audio.Analyze(path, audio.Options{Language: mcp.ParseString(req, "language", "auto")})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return resultJSON(res)
}
