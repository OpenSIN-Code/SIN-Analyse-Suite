package main

import (
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sin-analyse",
	Short: "Multimodal preprocessing pipelines for SIN-Code",
	Long:  "SIN-Analyse-Suite provides analysis pipelines for images, videos, PDFs, logs, data files, and audio.",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(imageCmd)
	rootCmd.AddCommand(videoCmd)
	rootCmd.AddCommand(pdfCmd)
	rootCmd.AddCommand(logsCmd)
	rootCmd.AddCommand(dataCmd)
	rootCmd.AddCommand(audioCmd)
	rootCmd.AddCommand(serveCmd)
}
