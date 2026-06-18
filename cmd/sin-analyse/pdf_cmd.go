package main

import (
	"os"
	"encoding/json"
	"fmt"
	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/pdf"
	"github.com/spf13/cobra"
)

var pdfCmd = &cobra.Command{
	Use:   "pdf [path]",
	Short: "Analyse a PDF file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		ocr, _ := cmd.Flags().GetBool("ocr")
		j, _ := cmd.Flags().GetBool("json")
		res, err := pdf.Analyze(path, pdf.Options{OCR: ocr})
		if err != nil {
			return err
		}
		if j {
			return json.NewEncoder(os.Stdout).Encode(res)
		}
		fmt.Println(res.ToMarkdown())
		return nil
	},
}

func init() {
	pdfCmd.Flags().Bool("ocr", false, "Run OCR for scanned pages")
	pdfCmd.Flags().Bool("json", false, "Output JSON")
}
