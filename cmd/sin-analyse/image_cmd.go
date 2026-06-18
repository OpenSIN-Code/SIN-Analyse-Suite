package main

import (
	"os"
	"encoding/json"
	"fmt"
	"github.com/OpenSIN-Code/SIN-Analyse-Suite/cmd/sin-analyse/internal/image"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image [path]",
	Short: "Analyse an image file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		ocr, _ := cmd.Flags().GetBool("ocr")
		j, _ := cmd.Flags().GetBool("json")
		res, err := image.Analyze(path, image.Options{OCR: ocr})
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
	imageCmd.Flags().Bool("ocr", false, "Run OCR via Tesseract")
	imageCmd.Flags().Bool("json", false, "Output JSON")
}
